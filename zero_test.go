// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"github.com/juju/tc"
)

func (s *CheckersS) TestIsZero(c *tc.C) {
	testInfo(c, tc.IsZero, "IsZero", []string{"obtained"})
	testCheck(c, tc.IsZero, true, "", nil)
	testCheck(c, tc.IsZero, true, "", int32(0))
	testCheck(c, tc.IsZero, false, "", int32(1))
	testCheck(c, tc.IsZero, true, "", "")
	testCheck(c, tc.IsZero, true, "", struct{}{})
	testCheck(c, tc.IsZero, true, "", (*struct{})(nil))
	testCheck(c, tc.IsZero, false, "", "a")
	testCheck(c, tc.IsZero, false, "", []int{})
	testCheck(c, tc.IsZero, true, "", []int(nil))
}
