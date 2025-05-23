// Copyright 2022 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"time"

	. "github.com/juju/tc"
)

type TimeSuite struct{}

var _ = InternalSuite(&TimeSuite{})

func (s *TimeSuite) TestBefore(c *C) {
	now := time.Now()
	c.Assert(now, Before, now.Add(time.Second))
	c.Assert(now, Not(Before), now.Add(-time.Second))

	result, msg := Before.Check([]any{time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Check(msg, Equals, `expected 2 parameters, received 1`)

	result, msg = Before.Check([]any{42, time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Equals, `obtained param: expected type time.Time, received type int`)

	result, msg = Before.Check([]any{time.Time{}, "wow"}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Matches, `want param: expected type time.Time, received type string`)
}

func (s *TimeSuite) TestAfter(c *C) {
	now := time.Now()
	c.Assert(now, Not(After), now.Add(time.Second))
	c.Assert(now, After, now.Add(-time.Second))

	result, msg := After.Check([]any{time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Check(msg, Equals, `expected 2 parameters, received 1`)

	result, msg = After.Check([]any{42, time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Equals, `obtained param: expected type time.Time, received type int`)

	result, msg = After.Check([]any{time.Time{}, "wow"}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Matches, `want param: expected type time.Time, received type string`)
}

func (s *TimeSuite) TestAlmost(c *C) {
	now := time.Now()
	c.Assert(now, Not(Almost), now.Add(1001*time.Millisecond))
	c.Assert(now, Almost, now.Add(-time.Second))
	c.Assert(now, Almost, now.Add(time.Second))

	result, msg := Almost.Check([]any{time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Check(msg, Equals, `expected 2 parameters, received 1`)

	result, msg = Almost.Check([]any{42, time.Time{}}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Equals, `obtained param: expected type time.Time, received type int`)

	result, msg = Almost.Check([]any{time.Time{}, "wow"}, nil)
	c.Assert(result, Equals, false)
	c.Assert(msg, Matches, `want param: expected type time.Time, received type string`)
}
