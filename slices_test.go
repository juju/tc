// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	. "github.com/juju/tc"
)

type orderedSuite struct{}

var _ = InternalSuite(&orderedSuite{})

func (s *orderedSuite) TestSame(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 1, 2, 3,
	}
	c.Assert(left, OrderedLeft[[]int](Equals), right)
	c.Assert(left, OrderedRight[[]int](Equals), right)
	c.Assert(left, OrderedMatch[[]int](Equals), right)
}

func (s *orderedSuite) TestSparse(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 6, 1, 6, 2, 6, 6, 6, 3, 6,
	}
	c.Assert(left, OrderedLeft[[]int](Equals), right)
	c.Assert(left, Not(OrderedMatch[[]int](Equals)), right)
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

func (s *orderedSuite) TestMissingOneLast(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 1, 2, 3, 4,
	}
	c.Assert(left, Not(OrderedMatch[[]int](Equals)), right)
	c.Assert(right, Not(OrderedMatch[[]int](Equals)), left)
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
	c.Assert(left, OrderedMatch[[]string](Matches), right)
}

type unorderedSuite struct{}

var _ = InternalSuite(&unorderedSuite{})

func (s *unorderedSuite) TestSame(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		0, 1, 2, 3,
	}
	c.Assert(left, UnorderedMatch[[]int](Equals), right)
}

func (s *unorderedSuite) TestSameDisordered(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		2, 3, 0, 1,
	}
	c.Assert(left, UnorderedMatch[[]int](Equals), right)
}

func (s *unorderedSuite) TestDisorderedAlmostMatch(c *C) {
	left := []int{
		0, 1, 2, 3, 1,
	}
	right := []int{
		2, 3, 0, 1, 0,
	}
	c.Assert(left, Not(UnorderedMatch[[]int](Equals)), right)
}

func (s *unorderedSuite) TestDisorderedAlmostMatchShort(c *C) {
	left := []int{
		0, 1, 2, 3,
	}
	right := []int{
		2, 3, 0, 1, 0,
	}
	c.Assert(left, Not(UnorderedMatch[[]int](Equals)), right)
}
