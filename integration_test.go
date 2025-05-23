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

// -----------------------------------------------------------------------
// Integration test suite.

type integrationS struct{}

var _ = InternalSuite(&integrationS{})

type integrationTestHelper struct{}

func (s *integrationTestHelper) TestMultiLineStringEqualFails(c *C) {
	c.Check("foo\nbar\nbaz\nboom\n", Equals, "foo\nbaar\nbaz\nboom\n")
}

func (s *integrationTestHelper) TestStringEqualFails(c *C) {
	c.Check("foo", Equals, "bar")
}

func (s *integrationTestHelper) TestIntEqualFails(c *C) {
	c.Check(42, Equals, 43)
}

type complexStruct struct {
	r, i int
}

func (s *integrationTestHelper) TestStructEqualFails(c *C) {
	c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})
}

func (s *integrationS) TestCountSuite(c *C) {
	suitesRun += 1
}

func (s *integrationS) TestOutput(c *C) {
	exitCode, output := runHelperSuite("integrationTestHelper")
	c.Check(exitCode, Equals, 1)

	c.Check(output.Status("TestIntEqualFails"), Equals, "FAIL")
	c.Check(output.Logs("TestIntEqualFails"), Equals,
		`    integration_test.go:51: 
            c.Check(42, Equals, 43)
        ... obtained int = 42
        ... expected int = 43`)

	c.Check(output.Status("TestMultiLineStringEqualFails"), Equals, "FAIL")
	c.Check(output.Logs("TestMultiLineStringEqualFails"), Equals,
		`    integration_test.go:43: 
            c.Check("foo\nbar\nbaz\nboom\n", Equals, "foo\nbaar\nbaz\nboom\n")
        ... obtained string = "" +
        ...     "foo\n" +
        ...     "bar\n" +
        ...     "baz\n" +
        ...     "boom\n"
        ... expected string = "" +
        ...     "foo\n" +
        ...     "baar\n" +
        ...     "baz\n" +
        ...     "boom\n"
        ... String difference:
        ...     [1]: "bar" != "baar"`)

	c.Check(output.Status("TestStringEqualFails"), Equals, "FAIL")
	c.Check(output.Logs("TestStringEqualFails"), Equals,
		`    integration_test.go:47: 
            c.Check("foo", Equals, "bar")
        ... obtained string = "foo"
        ... expected string = "bar"`)

	c.Check(output.Status("TestStructEqualFails"), Equals, "FAIL")
	c.Check(output.Logs("TestStructEqualFails"), Equals,
		`    integration_test.go:59: 
            c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})
        ... obtained tc_test.complexStruct = tc_test.complexStruct{r:1, i:2}
        ... expected tc_test.complexStruct = tc_test.complexStruct{r:3, i:4}
        ... Difference:
        ...     r: 1 != 3
        ...     i: 2 != 4`)
}
