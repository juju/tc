// Copyright 2023 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"errors"
	"fmt"
	"reflect"
)

type errorIsChecker struct {
	*CheckerInfo
}

// ErrorIs checks whether a value is an error that matches the other
// argument.
var ErrorIs Checker = &errorIsChecker{
	CheckerInfo: &CheckerInfo{
		Name:   "ErrorIs",
		Params: []string{"obtained", "error"},
	},
}

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

func (checker *errorIsChecker) Check(params []any, names []string) (result bool, err string) {
	if params[1] == nil || params[0] == nil {
		return params[1] == params[0], ""
	}

	f := reflect.ValueOf(params[1])
	ft := f.Type()
	if !ft.Implements(errType) {
		return false, fmt.Sprintf("wrong error target type, got: %s", ft)
	}

	v := reflect.ValueOf(params[0])
	vt := v.Type()
	if !v.IsValid() {
		return false, fmt.Sprintf("wrong argument type %s for %s", vt, ft)
	}
	if !vt.Implements(errType) {
		return false, fmt.Sprintf("wrong argument type %s for %s", vt, ft)
	}

	return errors.Is(v.Interface().(error), f.Interface().(error)), ""
}
