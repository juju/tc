// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"github.com/juju/tc"
)

func (s *CheckersS) TestDeref(c *tc.C) {
	v := "str"
	testInfo(c, tc.Deref(tc.Equals), "Deref(obtained)=>Equals", []string{"obtained", "expected"})
	testCheck(c, tc.Equals, false, "Difference:\n...     *string != string", &v, v)
	testCheck(c, tc.Deref(tc.Equals), true, "", &v, v)
	testCheck(c, tc.Deref(tc.Equals), false, "obtained nil", nil, v)
}
