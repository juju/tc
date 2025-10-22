package tc_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"sync/atomic"

	. "github.com/juju/tc"
)

type MustSuite struct{}

var _ = InternalSuite(&MustSuite{})

func (s *MustSuite) TestMust(c *C) {
	Must(c, func() (string, error) {
		return "", nil
	})
}

func (s *MustSuite) TestMustFails(c *C) {
	pc := &panicC{c: c}
	func() {
		defer pc.recover()
		Must(pc, func() (string, error) {
			return "", errors.New("ouch")
		})
	}()
	c.Assert(pc.failed.Load(), IsTrue)
}

func (s *MustSuite) TestMust0(c *C) {
	Must0(c, func() (string, error) {
		return "", nil
	})
}

func (s *MustSuite) TestMust0Fails(c *C) {
	pc := &panicC{c: c}
	func() {
		defer pc.recover()
		Must0(pc, func() (string, error) {
			return "", errors.New("ouch")
		})
	}()
	c.Assert(pc.failed.Load(), IsTrue)
}

func (s *MustSuite) TestMust1(c *C) {
	r1 := Must1(c, func(a string) (string, error) {
		return a, nil
	}, "wow")
	c.Assert(r1, Equals, "wow")
}

func (s *MustSuite) TestMust1Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	func() {
		defer pc.recover()
		r1 = Must1(pc, func(a string) (string, error) {
			return a, errors.New("ouch")
		}, "wow")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
}

func (s *MustSuite) TestMust2(c *C) {
	r1 := Must2(c, func(a string, b string) (string, error) {
		return a + b, nil
	}, "wow", "cool")
	c.Assert(r1, Equals, "wowcool")
}

func (s *MustSuite) TestMust2Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	func() {
		defer pc.recover()
		r1 = Must2(pc, func(a string, b string) (string, error) {
			return a + b, errors.New("ouch")
		}, "wow", "cool")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
}

func (s *MustSuite) TestMust0_0(c *C) {
	Must0_0(c, func() error {
		return nil
	})
}

func (s *MustSuite) TestMust0_0Fails(c *C) {
	pc := &panicC{c: c}
	func() {
		defer pc.recover()
		Must0_0(pc, func() error {
			return errors.New("ouch")
		})
	}()
	c.Assert(pc.failed.Load(), IsTrue)
}

func (s *MustSuite) TestMust1_0(c *C) {
	var a1 string
	Must1_0(c, func(a string) error {
		a1 = a
		return nil
	}, "wow")
	c.Assert(a1, Equals, "wow")
}

func (s *MustSuite) TestMust1_0Fails(c *C) {
	pc := &panicC{c: c}
	var a1 string
	func() {
		defer pc.recover()
		Must1_0(pc, func(a string) error {
			a1 = a
			return errors.New("ouch")
		}, "wow")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(a1, Equals, "wow")
}

func (s *MustSuite) TestMust2_0(c *C) {
	var a1, a2 string
	Must2_0(c, func(a, b string) error {
		a1 = a
		a2 = b
		return nil
	}, "wow", "cool")
	c.Assert(a1, Equals, "wow")
	c.Assert(a2, Equals, "cool")
}

func (s *MustSuite) TestMust2_0Fails(c *C) {
	pc := &panicC{c: c}
	var a1, a2 string
	func() {
		defer pc.recover()
		Must2_0(pc, func(a, b string) error {
			a1 = a
			a2 = b
			return errors.New("ouch")
		}, "wow", "cool")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(a1, Equals, "wow")
	c.Assert(a2, Equals, "cool")
}

func (s *MustSuite) TestMust0_1(c *C) {
	Must0_1(c, func() (string, error) {
		return "", nil
	})
}

func (s *MustSuite) TestMust0_1Fails(c *C) {
	pc := &panicC{c: c}
	func() {
		defer pc.recover()
		Must0_1(pc, func() (string, error) {
			return "", errors.New("ouch")
		})
	}()
	c.Assert(pc.failed.Load(), IsTrue)
}

func (s *MustSuite) TestMust1_1(c *C) {
	r1 := Must1_1(c, func(a string) (string, error) {
		return a, nil
	}, "wow")
	c.Assert(r1, Equals, "wow")
}

func (s *MustSuite) TestMust1_1Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	func() {
		defer pc.recover()
		r1 = Must1_1(pc, func(a string) (string, error) {
			return a, errors.New("ouch")
		}, "wow")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
}

func (s *MustSuite) TestMust2_1(c *C) {
	r1 := Must2_1(c, func(a string, b string) (string, error) {
		return a + b, nil
	}, "wow", "cool")
	c.Assert(r1, Equals, "wowcool")
}

func (s *MustSuite) TestMust2_1Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	func() {
		defer pc.recover()
		r1 = Must2_1(pc, func(a string, b string) (string, error) {
			return a + b, errors.New("ouch")
		}, "wow", "cool")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
}

func (s *MustSuite) TestMust0_2(c *C) {
	Must0_2(c, func() (string, int, error) {
		return "", 0, nil
	})
}

func (s *MustSuite) TestMust0_2Fails(c *C) {
	pc := &panicC{c: c}
	func() {
		defer pc.recover()
		Must0_2(pc, func() (string, int, error) {
			return "", 0, errors.New("ouch")
		})
	}()
	c.Assert(pc.failed.Load(), IsTrue)
}

func (s *MustSuite) TestMust1_2(c *C) {
	r1, r2 := Must1_2(c, func(a string) (string, int, error) {
		return a, 1, nil
	}, "wow")
	c.Assert(r1, Equals, "wow")
	c.Assert(r2, Equals, 1)
}

func (s *MustSuite) TestMust1_2Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	var r2 int
	func() {
		defer pc.recover()
		r1, r2 = Must1_2(pc, func(a string) (string, int, error) {
			return a, 1, errors.New("ouch")
		}, "wow")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
	c.Assert(r2, Equals, 0)
}

func (s *MustSuite) TestMust2_2(c *C) {
	r1, r2 := Must2_2(c, func(a string, b string) (string, int, error) {
		return a + b, 1, nil
	}, "wow", "cool")
	c.Assert(r1, Equals, "wowcool")
	c.Assert(r2, Equals, 1)
}

func (s *MustSuite) TestMust2_2Fails(c *C) {
	pc := &panicC{c: c}
	var r1 string
	var r2 int
	func() {
		defer pc.recover()
		r1, r2 = Must2_2(pc, func(a string, b string) (string, int, error) {
			return a + b, 1, errors.New("ouch")
		}, "wow", "cool")
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(r1, Equals, "")
	c.Assert(r2, Equals, 0)
}

type panicC struct {
	c      LikeC
	failed atomic.Bool
	err    bytes.Buffer
}

func (pc *panicC) recover() {
	//revive:disable
	if err := recover(); err != nil {
		pc.failed.Store(true)
		if err != "" {
			fmt.Fprintln(&pc.err, err)
		}
	}
	//revive:enable
}

func (pc *panicC) Failed() bool {
	return pc.failed.Load()
}

func (pc *panicC) Fail() {
	pc.failed.Store(true)
}

func (pc *panicC) FailNow() {
	panic("")
}

func (pc *panicC) Fatal(args ...any) {
	panic(fmt.Sprint(args...))
}

func (pc *panicC) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func (pc *panicC) SkipNow() {
	pc.FailNow()
}

func (pc *panicC) Skipf(format string, args ...any) {
	pc.Fatalf(format, args...)
}

func (pc *panicC) Skipped() bool {
	return pc.failed.Load()
}

func (pc *panicC) Skip(args ...any) {
	pc.Fatal(args...)
}

func (pc *panicC) Assert(obtained any, checker Checker, args ...any) {
	params := append([]any{obtained}, args...)
	ok, errString := checker.Check(params, slices.Clone(checker.Info().Params))
	if !ok {
		panic(errString)
	}
}

func (pc *panicC) Check(obtained any, checker Checker, args ...any) bool {
	params := append([]any{obtained}, args...)
	ok, errString := checker.Check(params, slices.Clone(checker.Info().Params))
	if !ok {
		pc.failed.Store(true)
		fmt.Fprintln(&pc.err, errString)
	}

	return ok
}

func (pc *panicC) Attr(key, value string) {
	pc.c.Attr(key, value)
}

func (pc *panicC) Cleanup(f func()) {
	pc.c.Cleanup(f)
}

func (pc *panicC) Error(args ...any) {
	fmt.Fprint(&pc.err, args...)
}

func (pc *panicC) Errorf(format string, args ...any) {
	fmt.Fprintf(&pc.err, format, args...)
}

func (pc *panicC) Helper() {
	pc.c.Helper()
}

func (pc *panicC) Log(args ...any) {
	pc.c.Log(args...)
}

func (pc *panicC) Logf(format string, args ...any) {
	pc.c.Logf(format, args...)
}

func (pc *panicC) Name() string {
	return pc.c.Name()
}

func (pc *panicC) Setenv(key, value string) {
	pc.c.Setenv(key, value)
}

func (pc *panicC) Chdir(dir string) {
	pc.c.Chdir(dir)
}

func (pc *panicC) TempDir() string {
	return pc.c.TempDir()
}

func (pc *panicC) Context() context.Context {
	return pc.c.Context()
}

func (pc *panicC) Output() io.Writer {
	return pc.c.Output()
}

func (pc *panicC) TestName() string {
	return pc.c.TestName()
}

func (pc *panicC) Logger() Logger {
	return pc.c.Logger()
}

func (pc *panicC) MkDir() string {
	return pc.c.MkDir()
}
