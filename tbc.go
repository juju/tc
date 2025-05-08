// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import "testing"

// TBC wraps a testing.TB and implements LikeC.
type TBC struct {
	testing.TB
}

var _ LikeC = (*TBC)(nil)

func (tbc *TBC) TestName() string {
	return tbc.Name()
}

func (tbc *TBC) Output(calldepth int, s string) error {
	tbc.Logf("[LOG] X:XX.XXX %s", s)
	return nil
}

func (tbc *TBC) Check(obtained any, checker Checker, args ...any) bool {
	tbc.Helper()
	return Check(tbc, obtained, checker, args...)
}

func (tbc *TBC) Assert(obtained any, checker Checker, args ...any) {
	tbc.Helper()
	Assert(tbc, obtained, checker, args...)
}

func (tbc *TBC) MkDir() string {
	tbc.Helper()
	return tbc.TempDir()
}
