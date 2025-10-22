// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"context"
	"io"
)

// LikeTB is a copy of testing.TB without the private method.
type LikeTB interface {
	Attr(key, value string)
	Cleanup(func())
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Setenv(key, value string)
	Chdir(dir string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
	Context() context.Context
	Output() io.Writer
}

// LikeC offers almost the same as C but only has the
// common methods between testing.T and testing.B.
type LikeC interface {
	LikeTB
	TestName() string
	Logger() Logger
	Check(obtained any, checker Checker, args ...any) bool
	Assert(obtained any, checker Checker, args ...any)
	MkDir() string
}
