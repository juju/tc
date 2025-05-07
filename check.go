// Gocheck - A rich testing framework for Go
//
// Copyright (c) 2010-2013 Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//	http://labix.org/gocheck

package tc

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

const (
	fixtureKd = iota
	testKd
)

type funcKind int

const (
	succeededSt = iota
	failedSt
	skippedSt
	panickedSt
	fixturePanickedSt
	missedSt
)

type funcStatus uint32

// A method value can't reach its own Method structure.
type methodType struct {
	reflect.Value
	Info reflect.Method
}

func newMethod(receiver reflect.Value, i int) *methodType {
	return &methodType{receiver.Method(i), receiver.Type().Method(i)}
}

func (method *methodType) PC() uintptr {
	return method.Info.Func.Pointer()
}

func (method *methodType) suiteName() string {
	t := method.Info.Type.In(0)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func (method *methodType) String() string {
	return method.suiteName() + "." + method.Info.Name
}

func (m *methodType) Call(c *C) {
	if m == nil {
		return
	}

	c.method = m
	defer func() {
		c.method = nil
	}()

	c.Helper()

	switch {
	case m.Info.Type.In(1) == reflect.TypeOf((*C)(nil)) && m.Info.Type.NumIn() == 2:
		m.Value.Call([]reflect.Value{reflect.ValueOf(c)})
	case m.Info.Type.In(1) == reflect.TypeOf((*testing.T)(nil)) && m.Info.Type.NumIn() == 2:
		m.Value.Call([]reflect.Value{reflect.ValueOf(c.T)})
	default:
		c.Fatalf("bad signature for method %s: %T", m.Info.Name, m.Interface())
	}
}

func (method *methodType) matches(re *regexp.Regexp) bool {
	return (re.MatchString(method.Info.Name) ||
		re.MatchString(method.suiteName()) ||
		re.MatchString(method.String()))
}

type C struct {
	*testing.T

	method    *methodType
	reason    string
	mustFail  bool
	tempDir   *tempDir
	startTime time.Time
}

// -----------------------------------------------------------------------
// Some simple formatting helpers.

var initWD, initWDErr = os.Getwd()

func init() {
	if initWDErr == nil {
		initWD = strings.Replace(initWD, "\\", "/", -1) + "/"
	}
}

func nicePath(path string) string {
	if initWDErr == nil {
		if strings.HasPrefix(path, initWD) {
			return path[len(initWD):]
		}
	}
	return path
}

func niceFuncPath(pc uintptr) string {
	function := runtime.FuncForPC(pc)
	if function != nil {
		filename, line := function.FileLine(pc)
		return fmt.Sprintf("%s:%d", nicePath(filename), line)
	}
	return "<unknown path>"
}

func niceFuncName(pc uintptr) string {
	function := runtime.FuncForPC(pc)
	if function != nil {
		name := path.Base(function.Name())
		if i := strings.Index(name, "."); i > 0 {
			name = name[i+1:]
		}
		if strings.HasPrefix(name, "(*") {
			if i := strings.Index(name, ")"); i > 0 {
				name = name[2:i] + name[i+1:]
			}
		}
		if i := strings.LastIndex(name, ".*"); i != -1 {
			name = name[:i] + "." + name[i+2:]
		}
		if i := strings.LastIndex(name, "·"); i != -1 {
			name = name[:i] + "." + name[i+2:]
		}
		return name
	}
	return "<unknown function>"
}

func suiteName(suite any) string {
	suiteType := reflect.TypeOf(suite)
	if suiteType.Kind() == reflect.Ptr {
		return suiteType.Elem().Name()
	}
	return suiteType.Name()
}

// -----------------------------------------------------------------------
// The underlying suite runner.

type suiteRunner struct {
	suite                     any
	setUpSuite, tearDownSuite *methodType
	setUpTest, tearDownTest   *methodType
	tests                     []*methodType
	tempDir                   *tempDir
	keepDir                   bool
}

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite any) *suiteRunner {
	suiteType := reflect.TypeOf(suite)
	suiteNumMethods := suiteType.NumMethod()
	suiteValue := reflect.ValueOf(suite)

	runner := &suiteRunner{
		suite:   suite,
		tempDir: &tempDir{},
		tests:   make([]*methodType, 0, suiteNumMethods),
	}

	for i := 0; i != suiteNumMethods; i++ {
		method := newMethod(suiteValue, i)
		switch method.Info.Name {
		case "SetUpSuite":
			runner.setUpSuite = method
		case "TearDownSuite":
			runner.tearDownSuite = method
		case "SetUpTest":
			runner.setUpTest = method
		case "TearDownTest":
			runner.tearDownTest = method
		default:
			if !strings.HasPrefix(method.Info.Name, "Test") {
				continue
			}
			runner.tests = append(runner.tests, method)
		}
	}
	return runner
}

// Run all methods in the given suite.
func (runner *suiteRunner) run(t *testing.T) {
	t.Cleanup(func() {
		runner.tempDir.removeAll()
	})

	c := C{T: t, startTime: time.Now()}
	t.Cleanup(func() { runner.tearDownSuite.Call(&c) })
	runner.setUpSuite.Call(&c)

	for _, test := range runner.tests {
		t.Run(test.Info.Name, func(t *testing.T) {
			runner.runTest(t, test)
		})
	}
}

// Same as forkTest(), but wait for the test to finish before returning.
func (runner *suiteRunner) runTest(t *testing.T, method *methodType) {
	c := C{T: t, startTime: time.Now()}

	t.Cleanup(func() { runner.tearDownTest.Call(&c) })
	runner.setUpTest.Call(&c)
	method.Call(&c)
}
