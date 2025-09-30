// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// TBC wraps a testing.TB and implements LikeC.
type TBC struct {
	testing.TB
}

var _ LikeC = (*TBC)(nil)

func (tbc *TBC) TestName() string {
	return tbc.Name()
}

func (tbc *TBC) Logger() Logger {
	return &tbcLogger{tbc}
}

type tbcLogger struct {
	*TBC
}

func (l *tbcLogger) Output(calldepth int, s string) error {
	_, file, line, _ := runtime.Caller(calldepth)
	file = filepath.Base(file)
	fmt.Fprintf(l.TB.Output(), "%s:%d: %s\n", file, line, s)
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
