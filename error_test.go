package check_test

import (
	"fmt"

	"github.com/juju/errors"

	. "gopkg.in/check.v2"
)

type ErrorSuite struct{}

var _ = Suite(&ErrorSuite{})

var errorIsTests = []struct {
	arg    interface{}
	target interface{}
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
	arg:    errors.ConstError("bar"),
	target: errors.ConstError("foo"),
	result: false,
}, {
	arg:    errors.ConstError("foo"),
	target: errors.ConstError("foo"),
	result: true,
}, {
	arg:    errors.Trace(errors.ConstError("foo")),
	target: errors.ConstError("foo"),
	result: true,
}, {
	arg:    errors.ConstError("foo"),
	target: "blah",
	msg:    "wrong error target type, got: string",
}, {
	arg:    "blah",
	target: errors.ConstError("foo"),
	msg:    "wrong argument type string for errors.ConstError",
}, {
	arg:    (*error)(nil),
	target: errors.ConstError("foo"),
	msg:    "wrong argument type *error for errors.ConstError",
}}

func (s *ErrorSuite) TestErrorIs(c *C) {
	for i, test := range errorIsTests {
		c.Logf("test %d. %T %T", i, test.arg, test.target)
		result, msg := ErrorIs.Check([]interface{}{test.arg, test.target}, nil)
		c.Check(result, Equals, test.result)
		c.Check(msg, Equals, test.msg)
	}
}
