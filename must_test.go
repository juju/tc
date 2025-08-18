package tc_test

import (
	"errors"
	"fmt"
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	pc := &panicC{LikeC: c}
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
	LikeC
	failed    atomic.Bool
	errString string
}

func (pc *panicC) recover() {
	//revive:disable
	if err := recover(); err != nil {
		if !pc.failed.Swap(true) {
			pc.errString = fmt.Sprintf("%v", err)
		}
	}
	//revive:enable
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
		if !pc.failed.Swap(true) {
			pc.errString = errString
		}
	}
	return ok
}
