// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"slices"
	"strings"
)

// Not checks that the provided checker does not pass.
func Not(checker Checker) Checker {
	return &logicalChecker{
		name:     "Not",
		checkers: []Checker{checker},
		op: func(outcomes []bool, _ []string) (bool, string) {
			return !outcomes[0], ""
		},
	}
}

// And checks that all of the provided checkers pass.
func And(checker Checker, checkers ...Checker) Checker {
	return &logicalChecker{
		name:     "And",
		checkers: append([]Checker{checker}, checkers...),
		op: func(outcomes []bool, errors []string) (bool, string) {
			res := true
			resErrors := []string{}
			for i, b := range outcomes {
				if !b && errors[i] != "" {
					resErrors = append(resErrors, errors[i])
				}
				res = res && b
			}
			return res, strings.Join(resErrors, "\n")
		},
	}
}

// Or checks that one of the provided checkers pass.
func Or(checker Checker, checkers ...Checker) Checker {
	return &logicalChecker{
		name:     "Or",
		checkers: append([]Checker{checker}, checkers...),
		op: func(outcomes []bool, errors []string) (bool, string) {
			res := false
			resErrors := []string{}
			for i, b := range outcomes {
				if !b && errors[i] != "" {
					resErrors = append(resErrors, errors[i])
				}
				res = res || b
			}
			if res {
				return true, ""
			}
			return false, strings.Join(resErrors, "\n")
		},
	}
}

type logicalChecker struct {
	name     string
	checkers []Checker
	op       func([]bool, []string) (bool, string)
}

func (c *logicalChecker) Info() *CheckerInfo {
	info := &CheckerInfo{}
	names := []string{}
	for _, v := range c.checkers {
		childInfo := v.Info()
		if len(info.Params) < len(childInfo.Params) {
			info.Params = slices.Clone(childInfo.Params)
		}
		names = append(names, childInfo.Name)
	}
	info.Name = fmt.Sprintf("%s(%s)", c.name, strings.Join(names, ", "))
	return info
}

func (c *logicalChecker) Check(params []any, names []string) (bool, string) {
	outcomes := make([]bool, len(c.checkers))
	errors := make([]string, len(c.checkers))
	for i, checker := range c.checkers {
		info := checker.Info()
		checkerParams := params[:len(info.Params)]
		outcomes[i], errors[i] = checker.Check(
			checkerParams, slices.Clone(info.Params))
	}
	return c.op(outcomes, errors)
}
