// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"github.com/juju/tc"
)

func (s *CheckersS) TestBindChecker(c *tc.C) {
	testInfo(c, tc.Bind(tc.Equals, "foo"), "Equals(#1)", []string{"obtained"})
	testCheck(c, tc.Bind(tc.Equals, "foo"), false, "", "bar")
	testCheck(c, tc.Bind(tc.Equals, "foo"), true, "", "foo")
	testCheck(c, tc.Bind(tc.DeepEquals, map[string]string{"foo": "bar"}), true, "", map[string]string{"foo": "bar"})
}

func (s *CheckersS) TestBindMatcher(c *tc.C) {
	bound := tc.Bind(tc.Equals, "foo")
	c.Check(bound.Matches("foo"), tc.IsTrue)
	c.Check(bound.Got("bar"), tc.Equals, "expected Equals(foo) got bar")
	c.Check(bound.Matches("bar"), tc.IsFalse)
	c.Check(bound.Matches(nil), tc.IsFalse)
	c.Check(bound.String(), tc.Equals, "Equals(foo)")
}
