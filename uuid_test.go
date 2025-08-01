// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"github.com/juju/tc"
)

func (s *CheckersS) TestIsUUID(c *tc.C) {
	testInfo(c, tc.IsUUID, "IsUUID", []string{"obtained"})
	testCheck(c, tc.IsUUID, false, "obtained value does not look like a uuid", "")
	testCheck(c, tc.IsUUID, true, "", "00000000-0000-0000-0000-000000000000")
	testCheck(c, tc.IsUUID, true, "", "35101a9e-1f8a-4e92-903b-a0616b931b79")
	testCheck(c, tc.IsUUID, false, "obtained value does not look like a uuid", ":0000000-0000-0000-0000-000000000000")
}

func (s *CheckersS) TestIsZeroUUID(c *tc.C) {
	testInfo(c, tc.IsZeroUUID, "IsZeroUUID", []string{"obtained"})
	testCheck(c, tc.IsZeroUUID, false, "obtained value does not look like a uuid", "")
	testCheck(c, tc.IsZeroUUID, true, "", "00000000-0000-0000-0000-000000000000")
	testCheck(c, tc.IsZeroUUID, false, "", "35101a9e-1f8a-4e92-903b-a0616b931b79")
	testCheck(c, tc.IsZeroUUID, false, "obtained value does not look like a uuid", ":0000000-0000-0000-0000-000000000000")
}

func (s *CheckersS) TestIsNonZeroUUID(c *tc.C) {
	testInfo(c, tc.IsNonZeroUUID, "And(Not(IsZeroUUID), IsUUID)", []string{"obtained"})
	testCheck(c, tc.IsNonZeroUUID, false, "obtained value does not look like a uuid", "")
	testCheck(c, tc.IsNonZeroUUID, false, "", "00000000-0000-0000-0000-000000000000")
	testCheck(c, tc.IsNonZeroUUID, true, "", "35101a9e-1f8a-4e92-903b-a0616b931b79")
	testCheck(c, tc.IsNonZeroUUID, false, "obtained value does not look like a uuid", ":0000000-0000-0000-0000-000000000000")
}
