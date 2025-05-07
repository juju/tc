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

package tc_test

import (
	. "github.com/juju/tc"
)

var runnerS = Suite(&RunS{})

type RunS struct{}

func (s *RunS) TestCountSuite(c *C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *RunS) TestSuccess(c *C) {
	exitCode, output := runHelperSuite("SuccessHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("TestLogAndSucceed"), Equals, "PASS")
}

func (s *RunS) TestFailure(c *C) {
	exitCode, output := runHelperSuite("FailHelper")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("TestLogAndFail"), Equals, "FAIL")
}

func (s *RunS) TestFixture(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("Test1"), Equals, "PASS")
	c.Check(output.Status("Test2"), Equals, "PASS")
}

func (s *RunS) TestPanicOnTest(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "Test1")
	c.Check(exitCode, Equals, 2)
	c.Check(output.Status("Test1"), Equals, "FAIL")
	// stdlib testing stops on first panic
	c.Check(output.Status("Test2"), Equals, "")
}

func (s *RunS) TestPanicOnSetUpTest(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "SetUpTest")
	c.Check(exitCode, Equals, 2)
	c.Check(output.Status("Test1"), Equals, "FAIL")
	// stdlib testing stops on first panic
	c.Check(output.Status("Test2"), Equals, "")
}

func (s *RunS) TestPanicOnSetUpSuite(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "SetUpSuite")
	c.Check(exitCode, Equals, 2)
	// If SetUpSuite fails, no tests from the suite are run
	c.Check(output.Status("Test1"), Equals, "")
	c.Check(output.Status("Test2"), Equals, "")
}

/*
// -----------------------------------------------------------------------
// Verify that List works correctly.

func (s *RunS) TestListFiltered(c *C) {
	names := List(&FixtureHelper{}, &RunConf{Filter: "1"})
	c.Assert(names, DeepEquals, []string{
		"FixtureHelper.Test1",
	})
}

func (s *RunS) TestList(c *C) {
	names := List(&FixtureHelper{}, &RunConf{})
	c.Assert(names, DeepEquals, []string{
		"FixtureHelper.Test1",
		"FixtureHelper.Test2",
	})
}

// -----------------------------------------------------------------------
// Verify that that the keep work dir request indeed does so.

type WorkDirSuite struct {}

func (s *WorkDirSuite) Test(c *C) {
	c.MkDir()
}

func (s *RunS) TestKeepWorkDir(c *C) {
	output := String{}
	runConf := RunConf{Output: &output, Verbose: true, KeepWorkDir: true}
	result := Run(&WorkDirSuite{}, &runConf)

	c.Assert(result.String(), Matches, ".*\nWORK=" + regexp.QuoteMeta(result.WorkDir))

	stat, err := os.Stat(result.WorkDir)
	c.Assert(err, IsNil)
	c.Assert(stat.IsDir(), Equals, true)
}
*/
