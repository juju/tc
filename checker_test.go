// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package check_test

import (
	"fmt"
	"time"

	. "gopkg.in/check.v2"
)

type CheckerSuite struct{}

var _ = Suite(&CheckerSuite{})

func (s *CheckerSuite) TestHasPrefix(c *C) {
	c.Assert("foo bar", HasPrefix, "foo")
	c.Assert("foo bar", Not(HasPrefix), "omg")
}

func (s *CheckerSuite) TestHasSuffix(c *C) {
	c.Assert("foo bar", HasSuffix, "bar")
	c.Assert("foo bar", Not(HasSuffix), "omg")
}

func (s *CheckerSuite) TestContains(c *C) {
	c.Assert("foo bar baz", Contains, "foo")
	c.Assert("foo bar baz", Contains, "bar")
	c.Assert("foo bar baz", Contains, "baz")
	c.Assert("foo bar baz", Not(Contains), "omg")
}

func (s *CheckerSuite) TestTimeBetween(c *C) {
	now := time.Now()
	earlier := now.Add(-1 * time.Second)
	later := now.Add(time.Second)

	checkOK := func(value interface{}, start, end time.Time) {
		checker := TimeBetween(start, end)
		value, msg := checker.Check([]interface{}{value}, nil)
		c.Check(value, IsTrue)
		c.Check(msg, Equals, "")
	}

	checkFails := func(value interface{}, start, end time.Time, match string) {
		checker := TimeBetween(start, end)
		value, msg := checker.Check([]interface{}{value}, nil)
		c.Check(value, IsFalse)
		c.Check(msg, Matches, match)
	}

	checkOK(now, earlier, later)
	// Later can be before earlier...
	checkOK(now, later, earlier)
	// check at bounds
	checkOK(earlier, earlier, later)
	checkOK(later, earlier, later)

	checkFails(earlier, now, later, `obtained time .* is before start time .*`)
	checkFails(later, now, earlier, `obtained time .* is after end time .*`)
	checkFails(42, now, earlier, `obtained value type must be time.Time`)
}

type someStruct struct {
	a uint
}

func (s *CheckerSuite) TestSameContents(c *C) {
	//// positive cases ////

	// same
	c.Check(
		[]int{1, 2, 3}, SameContents,
		[]int{1, 2, 3})

	// empty
	c.Check(
		[]int{}, SameContents,
		[]int{})

	// single
	c.Check(
		[]int{1}, SameContents,
		[]int{1})

	// different order
	c.Check(
		[]int{1, 2, 3}, SameContents,
		[]int{3, 2, 1})

	// multiple copies of same
	c.Check(
		[]int{1, 1, 2}, SameContents,
		[]int{2, 1, 1})

	type test struct {
		s string
		i int
	}

	// test structs
	c.Check(
		[]test{{"a", 1}, {"b", 2}}, SameContents,
		[]test{{"b", 2}, {"a", 1}})

	//// negative cases ////

	// different contents
	c.Check(
		[]int{1, 3, 2, 5}, Not(SameContents),
		[]int{5, 2, 3, 4})

	// different size slices
	c.Check(
		[]int{1, 2, 3}, Not(SameContents),
		[]int{1, 2})

	// different counts of same items
	c.Check(
		[]int{1, 1, 2}, Not(SameContents),
		[]int{1, 2, 2})

	// Tests that check that we compare the contents of structs,
	// that we point to, not just the pointers to them.
	a1 := someStruct{1}
	a2 := someStruct{2}
	a3 := someStruct{3}
	b1 := someStruct{1}
	b2 := someStruct{2}
	// Same order, same contents
	c.Check(
		[]*someStruct{&a1, &a2}, SameContents,
		[]*someStruct{&b1, &b2})

	// Empty vs not
	c.Check(
		[]*someStruct{&a1, &a2}, Not(SameContents),
		[]*someStruct{})

	// Empty vs empty
	// Same order, same contents
	c.Check(
		[]*someStruct{}, SameContents,
		[]*someStruct{})

	// Different order, same contents
	c.Check(
		[]*someStruct{&a1, &a2}, SameContents,
		[]*someStruct{&b2, &b1})

	// different contents
	c.Check(
		[]*someStruct{&a3, &a2}, Not(SameContents),
		[]*someStruct{&b2, &b1})

	// Different sizes, same contents (duplicate item)
	c.Check(
		[]*someStruct{&a1, &a2, &a1}, Not(SameContents),
		[]*someStruct{&b2, &b1})

	// Different sizes, same contents
	c.Check(
		[]*someStruct{&a1, &a1, &a2}, Not(SameContents),
		[]*someStruct{&b2, &b1})

	// Same sizes, same contents, different quantities
	c.Check(
		[]*someStruct{&a1, &a2, &a2}, Not(SameContents),
		[]*someStruct{&b1, &b1, &b2})

	/// Error cases ///
	//  note: for these tests, we can't use Not, since Not passes the error value through
	// and checks with a non-empty error always count as failed
	// Oddly, there doesn't seem to actually be a way to check for an error from a Checker.

	// different type
	res, err := SameContents.Check([]interface{}{
		[]string{"1", "2"},
		[]int{1, 2},
	}, []string{})
	c.Check(res, IsFalse)
	c.Check(err, Not(Equals), "")

	// obtained not a slice
	res, err = SameContents.Check([]interface{}{
		"test",
		[]int{1},
	}, []string{})
	c.Check(res, IsFalse)
	c.Check(err, Not(Equals), "")

	// expected not a slice
	res, err = SameContents.Check([]interface{}{
		[]int{1},
		"test",
	}, []string{})
	c.Check(res, IsFalse)
	c.Check(err, Not(Equals), "")
}

type stack_error struct {
	message string
	stack   []string
}

type embedded struct {
	typed *stack_error
	err   error
}

func (s *stack_error) Error() string {
	return s.message
}
func (s *stack_error) StackTrace() []string {
	return s.stack
}

type value_error string

func (e value_error) Error() string {
	return string(e)
}

func (s *CheckerSuite) TestErrorIsNil(c *C) {
	checkOK := func(value interface{}) {
		value, msg := ErrorIsNil.Check([]interface{}{value}, nil)
		c.Check(value, IsTrue)
		c.Check(msg, Equals, "")
	}

	checkFails := func(value interface{}, match string) {
		value, msg := ErrorIsNil.Check([]interface{}{value}, nil)
		c.Check(value, IsFalse)
		c.Check(msg, Matches, match)
	}

	var typedNil *stack_error
	var typedNilAsInterface error = typedNil
	var nilError error
	var value value_error
	var emptyValueErrorAsInterface error = value
	var embed embedded

	checkOK(nil)
	checkOK(nilError)
	checkOK(embed.err)

	checkFails([]string{}, `obtained type \(.*\) is not an error`)
	checkFails("", `obtained type \(.*\) is not an error`)
	checkFails(embed.typed, `value of \(.*\) is nil, but a typed nil`)
	checkFails(typedNilAsInterface, `value of \(.*\) is nil, but a typed nil`)
	checkFails(fmt.Errorf("an error"), "")
	checkFails(value, "")
	checkFails(emptyValueErrorAsInterface, "")

	emptyStack := &stack_error{"message", nil}
	checkFails(emptyStack, "")

	withStack := &stack_error{"message", []string{
		"filename:line", "filename2:line2"}}
	checkFails(withStack, "error stack:\n\tfilename:line\n\tfilename2:line2")
}
