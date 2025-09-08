// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	. "github.com/juju/tc"
)

type MultiCheckerSuite struct{}

var _ = InternalSuite(&MultiCheckerSuite{})

func (s *MultiCheckerSuite) TestDeepEquals(c *C) {
	for i, test := range deepEqualTests {
		c.Logf("test %d. %v == %v is %v", i, test.a, test.b, test.eq)
		result, msg := NewMultiChecker().Check([]any{test.a, test.b}, nil)
		c.Check(result, Equals, test.eq)
		if test.eq {
			c.Check(msg, Equals, "")
		} else {
			c.Check(msg, Not(Equals), "")
		}
	}
}

func (s *MultiCheckerSuite) TestArray(c *C) {
	a1 := []string{"a", "b", "c"}
	a2 := []string{"a", "bbb", "c"}

	checker := NewMultiChecker().AddExpr("_[1]", Ignore)
	c.Check(a1, checker, a2)
}

func (s *MultiCheckerSuite) TestMap(c *C) {
	a1 := map[string]string{"a": "a", "b": "b", "c": "c"}
	a2 := map[string]string{"a": "a", "b": "bbbb", "c": "c"}

	checker := NewMultiChecker().AddExpr(`_["b"]`, Ignore)
	c.Check(a1, checker, a2)
}

func (s *MultiCheckerSuite) TestArrayArraysUnordered(c *C) {
	a1 := [][]string{{"a", "b", "c"}, {"c", "d", "e"}}
	a2 := [][]string{{"a", "b", "c"}, {}}

	checker := NewMultiChecker().AddExpr("_[1]", SameContents, []string{"e", "c", "d"})
	c.Check(a1, checker, a2)
}

func (s *MultiCheckerSuite) TestArrayArraysUnorderedWithExpected(c *C) {
	a1 := [][]string{{"a", "b", "c"}, {"c", "d", "e"}}
	a2 := [][]string{{"a", "b", "c"}, {"e", "c", "d"}}

	checker := NewMultiChecker().AddExpr("_[1]", SameContents, ExpectedValue)
	c.Check(a1, checker, a2)
}

type pod struct {
	A int
	a int
	B bool
	b bool
	C string
	c string
}

func (s *MultiCheckerSuite) TestPOD(c *C) {
	a1 := pod{1, 2, true, true, "a", "a"}
	a2 := pod{2, 3, false, false, "b", "b"}

	checker := NewMultiChecker().
		AddExpr("_.A", Ignore).
		AddExpr("_.a", Ignore).
		AddExpr("_.B", Ignore).
		AddExpr("_.b", Ignore).
		AddExpr("_.C", Ignore).
		AddExpr("_.c", Ignore)
	c.Check(a1, checker, a2)
}

func (s *MultiCheckerSuite) TestExprMap(c *C) {
	a1 := map[string]string{"a": "a", "b": "b", "c": "c"}
	a2 := map[string]string{"a": "aaaa", "b": "bbbb", "c": "cccc"}

	checker := NewMultiChecker().AddExpr(`_[_]`, Ignore)
	c.Check(a1, checker, a2)
}

type complexA struct {
	complexB
	A int
	C []int
	D map[string]string
	E *complexE
	F **complexF
}

type complexB struct {
	B string
	b string
}

type complexE struct {
	E string
}

type complexF struct {
	F []string
}

func (s *MultiCheckerSuite) TestExprComplex(c *C) {
	f1 := &complexF{
		F: []string{"a", "b"},
	}
	a1 := complexA{
		complexB: complexB{
			B: "wow",
			b: "wow",
		},
		A: 5,
		C: []int{0, 1, 2, 3, 4, 5},
		D: map[string]string{"a": "b"},
		E: &complexE{E: "E"},
		F: &f1,
	}
	f2 := &complexF{
		F: []string{"c", "d"},
	}
	a2 := complexA{
		complexB: complexB{
			B: "cool",
			b: "cool",
		},
		A: 19,
		C: []int{5, 4, 3, 2, 1, 0},
		D: map[string]string{"b": "a"},
		E: &complexE{E: "EEEEEEEEE"},
		F: &f2,
	}
	checker := NewMultiChecker().
		AddExpr(`_.complexB.B`, Ignore).
		AddExpr(`_.complexB.b`, Ignore).
		AddExpr(`_.A`, Ignore).
		AddExpr(`_.C[_]`, Ignore).
		AddExpr(`_.D`, Ignore).
		AddExpr(`(*_.E)`, Ignore).
		AddExpr(`(*(*_.F)).F[_]`, Ignore)
	c.Check(a1, checker, a2)
}

func (s *MultiCheckerSuite) TestExprComplexDefault(c *C) {
	f1 := &complexF{
		F: []string{"a", "b"},
	}
	a1 := complexA{
		complexB: complexB{
			B: "wow",
			b: "wow",
		},
		A: 5,
		C: []int{0, 1, 2, 3, 4, 5},
		D: map[string]string{
			"a":     "b",
			"VALID": "YES",
		},
		E: &complexE{E: "E"},
		F: &f1,
	}
	f2 := &complexF{
		F: []string{"c", "d"},
	}
	a2 := complexA{
		complexB: complexB{
			B: "cool",
			b: "cool",
		},
		A: 19,
		C: []int{5, 4, 3, 2, 1, 0},
		D: map[string]string{
			"b":     "a",
			"VALID": "YES",
		},
		E: &complexE{E: "EEEEEEEEE"},
		F: &f2,
	}
	checker := NewMultiChecker().
		SetDefault(Ignore).
		AddExpr(`_.D["VALID"]`, Equals, ExpectedValue)
	c.Check(a1, checker, a2)

	// Check it fails when inverting the passing check.
	checkerInverted := NewMultiChecker().
		SetDefault(Ignore).
		AddExpr(`_.D["VALID"]`, Not(Equals), ExpectedValue)
	pc := panicC{LikeC: c}
	func() {
		defer pc.recover()
		pc.Check(a1, checkerInverted, a2)
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(pc.errString, Contains, `.D["VALID"]`)
}

func (s *MultiCheckerSuite) TestExprLen(c *C) {
	a1 := []int{0, 1, 2, 3}
	a2 := []int{0, 1, 2, 3, 4}

	mc := NewMultiChecker()
	pc := panicC{LikeC: c}
	func() {
		defer pc.recover()
		pc.Check(a1, mc, a2)
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(pc.errString, Contains, "slice/array length mismatch")

	mcWithLenIgnore := NewMultiChecker().
		AddExpr(`len(_)`, Ignore)
	c.Assert(a1, mcWithLenIgnore, a2)
}

func (s *MultiCheckerSuite) TestMultipleMatches(c *C) {
	a1 := []any{1, 2, 3, 4.1}
	a2 := []any{1, 2, 3, 4.1}
	a3 := []any{1, 2, 3, 4.2}

	mc := NewMultiChecker().
		AddExpr(`_[_]`, GreaterThan, 0).
		AddExpr(`_[_]`, Equals, ExpectedValue).
		AddExpr(`_[3]`, FitsTypeOf, 0.0)
	c.Assert(a1, mc, a2)

	pc := panicC{LikeC: c}
	func() {
		defer pc.recover()
		pc.Check(a1, mc, a3)
	}()
	c.Assert(pc.failed.Load(), IsTrue)
	c.Assert(pc.errString, Equals, "mismatch at [3]: unequal; obtained 4.1; expected 4.2")
}
