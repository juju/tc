// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import "testing"

// LikeC offers almost the same as C but only has the
// common methods between testing.T and testing.B.
type LikeC interface {
	testing.TB
	TestName() string
	Output(calldepth int, s string) error
	Check(obtained any, checker Checker, args ...any) bool
	Assert(obtained any, checker Checker, args ...any)
	MkDir() string
}
