// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	. "github.com/juju/tc"
)

type RelopSuite struct{}

var _ = Suite(&RelopSuite{})

func (s *RelopSuite) TestGreaterThan(c *C) {
	c.Assert(45, GreaterThan, 42)
	c.Assert(2.25, GreaterThan, 1.0)
	c.Assert(42, Not(GreaterThan), 42)
	c.Assert(10, Not(GreaterThan), 42)

	result, msg := GreaterThan.Check([]any{"Hello", "World"}, nil)
	c.Assert(result, IsFalse)
	c.Assert(msg, Equals, `obtained value string:"Hello" not supported`)
}

func (s *RelopSuite) TestLessThan(c *C) {
	c.Assert(42, LessThan, 45)
	c.Assert(1.0, LessThan, 2.25)
	c.Assert(42, Not(LessThan), 42)
	c.Assert(42, Not(LessThan), 10)

	result, msg := LessThan.Check([]any{"Hello", "World"}, nil)
	c.Assert(result, IsFalse)
	c.Assert(msg, Equals, `obtained value string:"Hello" not supported`)
}
