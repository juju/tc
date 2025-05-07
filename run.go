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

package tc

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"testing"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []any

// Suite registers the given value as a test suite to be run. Any methods
// starting with the Test prefix in the given value will be considered as
// a test method.
func Suite(suite any) any {
	allSuites = append(allSuites, suite)
	return suite
}

// -----------------------------------------------------------------------
// Public running interface.

var (
	newListFlag = flag.Bool("tc.list", false, "List the names of all tests that will be run")
)

// TestingT runs all test suites registered with the Suite function,
// printing results to stdout, and reporting any failures back to
// the "testing" package.
func TestingT(t *testing.T) {
	t.Helper()
	if *newListFlag {
		w := bufio.NewWriter(os.Stdout)
		for _, name := range ListAll() {
			fmt.Fprintln(w, name)
		}
		w.Flush()
		return
	}
	RunAll(t)
}

// RunAll runs all test suites registered with the Suite function, using the
// provided run configuration.
func RunAll(t *testing.T) {
	t.Helper()
	for _, suite := range allSuites {
		t.Run(suiteName(suite), func(t *testing.T) {
			Run(t, suite)
		})
	}
}

// Run runs the provided test suite using the provided run configuration.
func Run(t *testing.T, suite any) {
	t.Helper()
	runner := newSuiteRunner(suite)
	runner.run(t)
}

// ListAll returns the names of all the test functions registered with the
// Suite function that will be run with the provided run configuration.
func ListAll() []string {
	var names []string
	for _, suite := range allSuites {
		names = append(names, List(suite)...)
	}
	return names
}

// List returns the names of the test functions in the given
// suite that will be run with the provided run configuration.
func List(suite any) []string {
	var names []string
	runner := newSuiteRunner(suite)
	for _, t := range runner.tests {
		names = append(names, t.String())
	}
	return names
}
