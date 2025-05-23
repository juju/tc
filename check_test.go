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

// This file contains just a few generic helpers which are used by the
// other test files.

package tc_test

import (
	"fmt"
	"os"

	"runtime"
	"testing"
	"time"

	"github.com/juju/tc"
)

// We count the number of suites run at least to get a vague hint that the
// test suite is behaving as it should.  Otherwise a bug introduced at the
// very core of the system could go unperceived.
const suitesRunExpected = 4

var suitesRun int = 0

func Test(t *testing.T) {
	tc.InternalTestingT(t)
	if suitesRun != suitesRunExpected {
		critical(fmt.Sprintf("Expected %d suites to run rather than %d",
			suitesRunExpected, suitesRun))
	}
}

// -----------------------------------------------------------------------
// Helper functions.

// Break down badly.  This is used in test cases which can't yet assume
// that the fundamental bits are working.
func critical(error string) {
	fmt.Fprintln(os.Stderr, "CRITICAL: "+error)
	os.Exit(1)
}

// Return the file line where it's called.
func getMyLine() int {
	if _, _, line, ok := runtime.Caller(1); ok {
		return line
	}
	return -1
}

// -----------------------------------------------------------------------
// Helper type implementing a basic io.Writer for testing output.

// Type implementing the io.Writer interface for analyzing output.
type String struct {
	value string
}

// The only function required by the io.Writer interface.  Will append
// written data to the String.value string.
func (s *String) Write(p []byte) (n int, err error) {
	s.value += string(p)
	return len(p), nil
}

// -----------------------------------------------------------------------
// Helper suite for testing basic fail behavior.

type FailHelper struct {
	testLine int
}

func (s *FailHelper) TestLogAndFail(c *tc.C) {
	s.testLine = getMyLine() - 1
	c.Log("Expected failure!")
	c.Fail()
}

// -----------------------------------------------------------------------
// Helper suite for testing basic success behavior.

type SuccessHelper struct{}

func (s *SuccessHelper) TestLogAndSucceed(c *tc.C) {
	c.Log("Expected success!")
}

// -----------------------------------------------------------------------
// Helper suite for testing ordering and behavior of fixture.

type FixtureHelper struct {
	calls   []string
	panicOn string
	skip    bool
	skipOnN int
	sleepOn string
	sleep   time.Duration
}

func (s *FixtureHelper) trace(name string, c *tc.C) {
	s.calls = append(s.calls, name)
	if name == s.panicOn {
		panic(name)
	}
	if s.sleep > 0 && s.sleepOn == name {
		time.Sleep(s.sleep)
	}
	if s.skip && s.skipOnN == len(s.calls)-1 {
		c.Skip("skipOnN == n")
	}
}

func (s *FixtureHelper) SetUpSuite(c *tc.C) {
	s.trace("SetUpSuite", c)
}

func (s *FixtureHelper) TearDownSuite(c *tc.C) {
	s.trace("TearDownSuite", c)
}

func (s *FixtureHelper) SetUpTest(c *tc.C) {
	s.trace("SetUpTest", c)
}

func (s *FixtureHelper) TearDownTest(c *tc.C) {
	s.trace("TearDownTest", c)
}

func (s *FixtureHelper) Test1(c *tc.C) {
	s.trace("Test1", c)
}

func (s *FixtureHelper) Test2(c *tc.C) {
	s.trace("Test2", c)
}

// -----------------------------------------------------------------------
// Helper which checks the state of the test and ensures that it matches
// the given expectations.  Depends on c.Errorf() working, so shouldn't
// be used to test this one function.

type SkippedSuite struct{}

var _ = tc.InternalSuite(&SkippedSuite{})

func (s *SkippedSuite) SetUpSuite(c *tc.C) {
	c.Skip("skippy")
}

func (s *SkippedSuite) TearDownSuite(c *tc.C) {
	c.FailNow()
}

func (s *SkippedSuite) SetUpTest(c *tc.C) {
	c.FailNow()
}

func (s *SkippedSuite) TearDownTest(c *tc.C) {
	c.FailNow()
}

func (s *SkippedSuite) TestShouldFail(c *tc.C) {
	c.FailNow()
}

type SkippedTestSuite struct{}

var _ = tc.InternalSuite(&SkippedTestSuite{})

func (s *SkippedTestSuite) SetUpTest(c *tc.C) {
	c.Skip("skippy")
}

func (s *SkippedTestSuite) TearDownTest(c *tc.C) {
	c.FailNow()
}

func (s *SkippedTestSuite) TestShouldFail(c *tc.C) {
	c.FailNow()
}
