// Copyright 2023 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc_test

import (
	"fmt"

	. "github.com/juju/tc"
)

// a ConstError is a prototype for a certain type of error
type ConstError string

// ConstError implements error
func (e ConstError) Error() string {
	return string(e)
}

type ErrorSuite struct{}

var _ = Suite(&ErrorSuite{})

var errorIsTests = []struct {
	arg    any
	target any
	result bool
	msg    string
}{{
	arg:    fmt.Errorf("bar"),
	target: nil,
	result: false,
}, {
	arg:    nil,
	target: fmt.Errorf("bar"),
	result: false,
}, {
	arg:    nil,
	target: nil,
	result: true,
}, {
	arg:    fmt.Errorf("bar"),
	target: fmt.Errorf("foo"),
	result: false,
}, {
	arg:    ConstError("bar"),
	target: ConstError("foo"),
	result: false,
}, {
	arg:    ConstError("foo"),
	target: ConstError("foo"),
	result: true,
}, {
	arg:    fmt.Errorf("%w", ConstError("foo")),
	target: ConstError("foo"),
	result: true,
}, {
	arg:    ConstError("foo"),
	target: "blah",
	msg:    "wrong error target type, got: string",
}, {
	arg:    "blah",
	target: ConstError("foo"),
	msg:    "wrong argument type string for tc_test.ConstError",
}, {
	arg:    (*error)(nil),
	target: ConstError("foo"),
	msg:    "wrong argument type *error for tc_test.ConstError",
}}

func (s *ErrorSuite) TestErrorIs(c *C) {
	for i, test := range errorIsTests {
		c.Logf("test %d. %T %T", i, test.arg, test.target)
		result, msg := ErrorIs.Check([]any{test.arg, test.target}, nil)
		c.Check(result, Equals, test.result)
		c.Check(msg, Equals, test.msg)
	}
}
