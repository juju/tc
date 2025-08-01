// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"github.com/juju/tc"
)

func (s *CheckersS) TestNot(c *tc.C) {
	testInfo(c, tc.Not(tc.IsTrue), "Not(IsTrue)", []string{"obtained"})
	testCheck(c, tc.Not(tc.IsFalse), false, "", false)
	testCheck(c, tc.Not(tc.IsFalse), true, "", true)
	testCheck(c, tc.Not(tc.Equals), true, "", 1, 2)
}

func (s *CheckersS) TestAnd(c *tc.C) {
	testInfo(c, tc.And(tc.IsTrue, tc.Equals), "And(IsTrue, Equals)", []string{"obtained", "expected"})
	testCheck(c, tc.And(tc.IsTrue, tc.Equals), false, "", false, false)
	testCheck(c, tc.And(tc.IsTrue, tc.Equals), true, "", true, true)
	testCheck(c, tc.And(tc.Equals), false, "", 1, 2)
}

func (s *CheckersS) TestOr(c *tc.C) {
	testInfo(c, tc.Or(tc.IsTrue, tc.Equals), "Or(IsTrue, Equals)", []string{"obtained", "expected"})
	testCheck(c, tc.Or(tc.IsTrue, tc.Equals), false, "", false, true)
	testCheck(c, tc.Or(tc.IsTrue, tc.Equals), true, "", true, true)
	testCheck(c, tc.Or(tc.Equals), false, "", 1, 2)
}
