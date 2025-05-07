// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"errors"
	"os"

	. "github.com/juju/tc"
)

type BoolSuite struct{}

var _ = Suite(&BoolSuite{})

func (s *BoolSuite) TestIsTrue(c *C) {
	c.Assert(true, IsTrue)
	c.Assert(false, Not(IsTrue))

	result, msg := IsTrue.Check([]any{false}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Equals, "")

	result, msg = IsTrue.Check([]any{"foo"}, nil)
	c.Assert(result, Equals, false)
	c.Check(msg, Equals, `expected type bool, received type string`)

	result, msg = IsTrue.Check([]any{42}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Equals, `expected type bool, received type int`)

	result, msg = IsTrue.Check([]any{nil}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Matches, `expected type bool, received <invalid .*Value>`)
}

func (s *BoolSuite) TestIsFalse(c *C) {
	c.Check(false, IsFalse)
	c.Check(true, Not(IsFalse))
}

func is42(i int) bool {
	return i == 42
}

var satisfiesTests = []struct {
	f      any
	arg    any
	result bool
	msg    string
}{{
	f:      is42,
	arg:    42,
	result: true,
}, {
	f:      is42,
	arg:    41,
	result: false,
}, {
	f:      is42,
	arg:    "",
	result: false,
	msg:    "wrong argument type string for func(int) bool",
}, {
	f:      os.IsNotExist,
	arg:    errors.New("foo"),
	result: false,
}, {
	f:      os.IsNotExist,
	arg:    os.ErrNotExist,
	result: true,
}, {
	f:      os.IsNotExist,
	arg:    nil,
	result: false,
}, {
	f:      func(chan int) bool { return true },
	arg:    nil,
	result: true,
}, {
	f:      func(func()) bool { return true },
	arg:    nil,
	result: true,
}, {
	f:      func(any) bool { return true },
	arg:    nil,
	result: true,
}, {
	f:      func(map[string]bool) bool { return true },
	arg:    nil,
	result: true,
}, {
	f:      func(*int) bool { return true },
	arg:    nil,
	result: true,
}, {
	f:      func([]string) bool { return true },
	arg:    nil,
	result: true,
}}

func (s *BoolSuite) TestSatisfies(c *C) {
	for i, test := range satisfiesTests {
		c.Logf("test %d. %T %T", i, test.f, test.arg)
		result, msg := Satisfies.Check([]any{test.arg, test.f}, nil)
		c.Check(result, Equals, test.result)
		c.Check(msg, Equals, test.msg)
	}
}

func (s *BoolSuite) TestDeepEquals(c *C) {
	for i, test := range deepEqualTests {
		c.Logf("test %d. %v == %v is %v", i, test.a, test.b, test.eq)
		result, msg := DeepEquals.Check([]any{test.a, test.b}, nil)
		c.Check(result, Equals, test.eq)
		if test.eq {
			c.Check(msg, Equals, "")
		} else {
			c.Check(msg, Not(Equals), "")
		}
	}
}
