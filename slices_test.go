// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	. "github.com/juju/tc"
)

type orderedSuite struct{}

var _ = Suite(&orderedSuite{})

func (s *orderedSuite) TestSame(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 1, 2, 3,
	}
	c.Assert(left, OrderedLeft[[]int](Equals), right)
}

func (s *orderedSuite) TestSparse(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 6, 1, 6, 2, 6, 6, 6, 3, 6,
	}
	c.Assert(left, OrderedLeft[[]int](Equals), right)
}

func (s *orderedSuite) TestMissing(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 6, 1, 6, 2, 6, 6, 6,
	}
	c.Assert(left, Not(OrderedLeft[[]int](Equals)), right)
}

func (s *orderedSuite) TestSubCheckerGetsExpectedValue(c *C) {
	left := []string{
		"apple", "orange", "bannana",
	}
	right := []string{
		"ap{2}le", "orange", "^ban{2}ana$",
	}
	c.Assert(left, OrderedLeft[[]string](Matches), right)
	c.Assert(left, OrderedRight[[]string](Matches), right)
}
