// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import "reflect"

// IsZero is a Checker that ensures the obtained value is the zero value for the
// obtained type.
var IsZero Checker = &zeroChecker{
	&CheckerInfo{Name: "IsZero", Params: []string{"obtained"}},
}

// NotZero is a Checker that ensures the obtained value is not the zero value
// for the obtained type.
var NotZero Checker = Not(IsZero)

type zeroChecker struct {
	*CheckerInfo
}

func (c *zeroChecker) Check(params []any, names []string) (bool, string) {
	obtained := params[0]
	value := reflect.ValueOf(obtained)
	if !value.IsValid() {
		// untyped nil is a zero any.
		return true, ""
	}
	return value.IsZero(), ""
}
