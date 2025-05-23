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
	"errors"
	"reflect"
	"runtime"

	"github.com/juju/tc"
)

type CheckersS struct{}

var _ = tc.InternalSuite(&CheckersS{})

func testInfo(c *tc.C, checker tc.Checker, name string, paramNames []string) {
	c.Helper()
	info := checker.Info()
	if info.Name != name {
		c.Fatalf("Got name %s, expected %s", info.Name, name)
	}
	if !reflect.DeepEqual(info.Params, paramNames) {
		c.Fatalf("Got param names %#v, expected %#v", info.Params, paramNames)
	}
}

func testCheck(c *tc.C, checker tc.Checker, result bool, error string, params ...any) ([]any, []string) {
	c.Helper()
	info := checker.Info()
	if len(params) != len(info.Params) {
		c.Fatalf("unexpected param count in test; expected %d got %d", len(info.Params), len(params))
	}
	names := append([]string{}, info.Params...)
	result_, error_ := checker.Check(params, names)
	if result_ != result || error_ != error {
		c.Fatalf("%s.Check(%#v) returned (%#v, %#v) rather than (%#v, %#v)",
			info.Name, params, result_, error_, result, error)
	}
	return params, names
}

func (s *CheckersS) TestCountSuite(c *tc.C) {
	suitesRun += 1
}

func (s *CheckersS) TestComment(c *tc.C) {
	bug := tc.Commentf("a %d bc", 42)
	comment := bug.CheckCommentString()
	if comment != "a 42 bc" {
		c.Fatalf("Commentf returned %#v", comment)
	}
}

func (s *CheckersS) TestIsNil(c *tc.C) {
	testInfo(c, tc.IsNil, "IsNil", []string{"value"})

	testCheck(c, tc.IsNil, true, "", nil)
	testCheck(c, tc.IsNil, false, "", "a")

	testCheck(c, tc.IsNil, true, "", (chan int)(nil))
	testCheck(c, tc.IsNil, false, "", make(chan int))
	testCheck(c, tc.IsNil, true, "", (error)(nil))
	testCheck(c, tc.IsNil, false, "", errors.New(""))
	testCheck(c, tc.IsNil, true, "", ([]int)(nil))
	testCheck(c, tc.IsNil, false, "", make([]int, 1))
	testCheck(c, tc.IsNil, false, "", int(0))
}

func (s *CheckersS) TestNotNil(c *tc.C) {
	testInfo(c, tc.NotNil, "NotNil", []string{"value"})

	testCheck(c, tc.NotNil, false, "", nil)
	testCheck(c, tc.NotNil, true, "", "a")

	testCheck(c, tc.NotNil, false, "", (chan int)(nil))
	testCheck(c, tc.NotNil, true, "", make(chan int))
	testCheck(c, tc.NotNil, false, "", (error)(nil))
	testCheck(c, tc.NotNil, true, "", errors.New(""))
	testCheck(c, tc.NotNil, false, "", ([]int)(nil))
	testCheck(c, tc.NotNil, true, "", make([]int, 1))
}

func (s *CheckersS) TestNot(c *tc.C) {
	testInfo(c, tc.Not(tc.IsNil), "Not(IsNil)", []string{"value"})

	testCheck(c, tc.Not(tc.IsNil), false, "", nil)
	testCheck(c, tc.Not(tc.IsNil), true, "", "a")
	testCheck(c, tc.Not(tc.Equals), true, "", 42, 43)
}

type simpleStruct struct {
	i int
}

func (s *CheckersS) TestEquals(c *tc.C) {
	testInfo(c, tc.Equals, "Equals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, tc.Equals, true, "", 42, 42)
	testCheck(c, tc.Equals, false, "", 42, 43)

	// Different native types.
	testCheck(c, tc.Equals, false, "", int32(42), int64(42))

	// With nil.
	testCheck(c, tc.Equals, false, "", 42, nil)
	testCheck(c, tc.Equals, false, "", nil, 42)
	testCheck(c, tc.Equals, true, "", nil, nil)

	// Slices
	testCheck(c, tc.Equals, false, "runtime error: comparing uncomparable type []uint8", []byte{1, 2}, []byte{1, 2})

	// Struct values
	testCheck(c, tc.Equals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, tc.Equals, false, `Difference:
...     i: 1 != 2`, simpleStruct{1}, simpleStruct{2})

	// Struct pointers, no difference in values, just pointer
	testCheck(c, tc.Equals, false, "", &simpleStruct{1}, &simpleStruct{1})
	// Struct pointers, different pointers and different values
	testCheck(c, tc.Equals, false, `Difference:
...     i: 1 != 2`, &simpleStruct{1}, &simpleStruct{2})
}

func (s *CheckersS) TestDeepEquals(c *tc.C) {
	testInfo(c, tc.DeepEquals, "DeepEquals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, tc.DeepEquals, true, "", 42, 42)
	testCheck(c, tc.DeepEquals, false, "mismatch at top level: unequal; obtained 42; expected 43", 42, 43)

	// Different native types.
	testCheck(c, tc.DeepEquals, false, "mismatch at top level: type mismatch int32 vs int64; obtained 42; expected 42", int32(42), int64(42))

	// With nil.
	testCheck(c, tc.DeepEquals, false, "mismatch at top level: nil vs non-nil mismatch; obtained 42; expected <nil>", 42, nil)

	// Slices
	testCheck(c, tc.DeepEquals, true, "", []byte{1, 2}, []byte{1, 2})
	testCheck(c, tc.DeepEquals, false, "mismatch at [1]: unequal; obtained 0x2; expected 0x3", []byte{1, 2}, []byte{1, 3})

	// Struct values
	testCheck(c, tc.DeepEquals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, tc.DeepEquals, false, "mismatch at .i: unequal; obtained 1; expected 2", simpleStruct{1}, simpleStruct{2})

	// Struct pointers
	testCheck(c, tc.DeepEquals, true, "", &simpleStruct{1}, &simpleStruct{1})
	s1 := &simpleStruct{1}
	s2 := &simpleStruct{2}
	testCheck(c, tc.DeepEquals, false, "mismatch at (*).i: unequal; obtained 1; expected 2", s1, s2)
}

func (s *CheckersS) TestHasLen(c *tc.C) {
	testInfo(c, tc.HasLen, "HasLen", []string{"obtained", "n"})

	testCheck(c, tc.HasLen, true, "", "abcd", 4)
	testCheck(c, tc.HasLen, true, "", []int{1, 2}, 2)
	testCheck(c, tc.HasLen, false, "", []int{1, 2}, 3)

	testCheck(c, tc.HasLen, false, "n must be an int", []int{1, 2}, "2")
	testCheck(c, tc.HasLen, false, "obtained value type has no length", nil, 2)
}

func (s *CheckersS) TestErrorMatches(c *tc.C) {
	testInfo(c, tc.ErrorMatches, "ErrorMatches", []string{"value", "regex"})

	testCheck(c, tc.ErrorMatches, false, "Error value is nil", nil, "some error")
	testCheck(c, tc.ErrorMatches, false, "Value is not an error", 1, "some error")
	testCheck(c, tc.ErrorMatches, true, "", errors.New("some error"), "some error")
	testCheck(c, tc.ErrorMatches, true, "", errors.New("some error"), "so.*or")

	// Verify params mutation
	params, names := testCheck(c, tc.ErrorMatches, false, "", errors.New("some error"), "other error")
	c.Assert(params[0], tc.Equals, "some error")
	c.Assert(names[0], tc.Equals, "error")
}

func (s *CheckersS) TestMatches(c *tc.C) {
	testInfo(c, tc.Matches, "Matches", []string{"value", "regex"})

	// Simple matching
	testCheck(c, tc.Matches, true, "", "abc", "abc")
	testCheck(c, tc.Matches, true, "", "abc", "a.c")

	// Must match fully
	testCheck(c, tc.Matches, false, "", "abc", "ab")
	testCheck(c, tc.Matches, false, "", "abc", "bc")

	// String()-enabled values accepted
	testCheck(c, tc.Matches, true, "", reflect.ValueOf("abc"), "a.c")
	testCheck(c, tc.Matches, false, "", reflect.ValueOf("abc"), "a.d")

	// Some error conditions.
	testCheck(c, tc.Matches, false, "Obtained value is not a string and has no .String()", 1, "a.c")
	testCheck(c, tc.Matches, false, "Can't compile regex: error parsing regexp: missing closing ]: `[c$`", "abc", "a[c")
}

func (s *CheckersS) TestPanics(c *tc.C) {
	testInfo(c, tc.Panics, "Panics", []string{"function", "expected"})

	// Some errors.
	testCheck(c, tc.Panics, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, tc.Panics, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, tc.Panics, true, "", func() { panic("BOOM") }, "BOOM")
	testCheck(c, tc.Panics, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, tc.Panics, true, "", func() bool { panic("BOOM") }, "BOOM")

	// Error values.
	testCheck(c, tc.Panics, true, "", func() { panic(errors.New("BOOM")) }, errors.New("BOOM"))
	testCheck(c, tc.Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))

	type deep struct{ i int }
	// Deep value
	testCheck(c, tc.Panics, true, "", func() { panic(&deep{99}) }, &deep{99})

	// Verify params/names mutation
	params, names := testCheck(c, tc.Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))
	c.Assert(params[0], tc.ErrorMatches, "KABOOM")
	c.Assert(names[0], tc.Equals, "panic")
}

func (s *CheckersS) TestPanicMatches(c *tc.C) {
	testInfo(c, tc.PanicMatches, "PanicMatches", []string{"function", "expected"})

	// Error matching.
	testCheck(c, tc.PanicMatches, true, "", func() { panic(errors.New("BOOM")) }, "BO.M")
	testCheck(c, tc.PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BO.M")

	// Some errors.
	testCheck(c, tc.PanicMatches, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, tc.PanicMatches, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, tc.PanicMatches, true, "", func() { panic("BOOM") }, "BO.M")
	testCheck(c, tc.PanicMatches, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, tc.PanicMatches, true, "", func() bool { panic("BOOM") }, "BO.M")

	// Verify params/names mutation
	params, names := testCheck(c, tc.PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BOOM")
	c.Assert(params[0], tc.Equals, "KABOOM")
	c.Assert(names[0], tc.Equals, "panic")
}

func (s *CheckersS) TestFitsTypeOf(c *tc.C) {
	testInfo(c, tc.FitsTypeOf, "FitsTypeOf", []string{"obtained", "sample"})

	// Basic types
	testCheck(c, tc.FitsTypeOf, true, "", 1, 0)
	testCheck(c, tc.FitsTypeOf, false, "", 1, int64(0))

	// Aliases
	testCheck(c, tc.FitsTypeOf, false, "", 1, errors.New(""))
	testCheck(c, tc.FitsTypeOf, false, "", "error", errors.New(""))
	testCheck(c, tc.FitsTypeOf, true, "", errors.New("error"), errors.New(""))

	// Structures
	testCheck(c, tc.FitsTypeOf, false, "", 1, simpleStruct{})
	testCheck(c, tc.FitsTypeOf, false, "", simpleStruct{42}, &simpleStruct{})
	testCheck(c, tc.FitsTypeOf, true, "", simpleStruct{42}, simpleStruct{})
	testCheck(c, tc.FitsTypeOf, true, "", &simpleStruct{42}, &simpleStruct{})

	// Some bad values
	testCheck(c, tc.FitsTypeOf, false, "Invalid sample value", 1, any(nil))
	testCheck(c, tc.FitsTypeOf, false, "", any(nil), 0)
}

func (s *CheckersS) TestImplements(c *tc.C) {
	testInfo(c, tc.Implements, "Implements", []string{"obtained", "ifaceptr"})

	var e error
	var re runtime.Error
	testCheck(c, tc.Implements, true, "", errors.New(""), &e)
	testCheck(c, tc.Implements, false, "", errors.New(""), &re)

	// Some bad values
	testCheck(c, tc.Implements, false, "ifaceptr should be a pointer to an interface variable", 0, errors.New(""))
	testCheck(c, tc.Implements, false, "ifaceptr should be a pointer to an interface variable", 0, any(nil))
	testCheck(c, tc.Implements, false, "", any(nil), &e)
}
