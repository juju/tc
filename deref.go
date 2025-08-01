// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"reflect"
	"slices"
)

// Deref dereferences the obtained value.
func Deref(checker Checker) Checker {
	return &derefChecker{
		checker: checker,
	}
}

type derefChecker struct {
	checker Checker
}

func (c *derefChecker) Info() *CheckerInfo {
	childInfo := c.checker.Info()
	info := &CheckerInfo{
		Name:   fmt.Sprintf("Deref(%s)=>%s", childInfo.Params[0], childInfo.Name),
		Params: slices.Clone(childInfo.Params),
	}
	return info
}

func (c *derefChecker) Check(params []any, names []string) (bool, string) {
	newParams := slices.Clone(params)
	obtained := newParams[0]
	if obtained == nil {
		return false, "obtained nil"
	}
	newParams[0] = reflect.Indirect(reflect.ValueOf(obtained)).Interface()
	return c.checker.Check(newParams, names)
}
