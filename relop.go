// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"reflect"
)

// GreaterThan checker

type greaterThanChecker struct {
	*CheckerInfo
}

var GreaterThan Checker = &greaterThanChecker{
	&CheckerInfo{Name: "GreaterThan", Params: []string{"obtained", "expected"}},
}

func (checker *greaterThanChecker) Check(params []any, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()

	gtZero := false
	if params[1] == 0 {
		gtZero = true
	}

	p0value := reflect.ValueOf(params[0])
	p1value := reflect.ValueOf(params[1])
	switch p0value.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if gtZero {
			return p0value.Int() > 0, ""
		}
		return p0value.Int() > p1value.Int(), ""
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if gtZero {
			return p0value.Uint() > 0, ""
		}
		return p0value.Uint() > p1value.Uint(), ""
	case reflect.Float32,
		reflect.Float64:
		if gtZero {
			return p0value.Float() > 0, ""
		}
		return p0value.Float() > p1value.Float(), ""
	default:
	}
	return false, fmt.Sprintf("obtained value %s:%#v not supported", p0value.Kind(), params[0])
}

// LessThan checker

type lessThanChecker struct {
	*CheckerInfo
}

var LessThan Checker = &lessThanChecker{
	&CheckerInfo{Name: "LessThan", Params: []string{"obtained", "expected"}},
}

func (checker *lessThanChecker) Check(params []any, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()

	ltZero := false
	if params[1] == 0 {
		ltZero = true
	}

	p0value := reflect.ValueOf(params[0])
	p1value := reflect.ValueOf(params[1])
	switch p0value.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if ltZero {
			return p0value.Int() < 0, ""
		}
		return p0value.Int() < p1value.Int(), ""
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if ltZero || p1value.Uint() == 0 {
			return false, "no possible value less than 0"
		}
		return p0value.Uint() < p1value.Uint(), ""
	case reflect.Float32,
		reflect.Float64:
		if ltZero {
			return p0value.Float() < 0, ""
		}
		return p0value.Float() < p1value.Float(), ""
	default:
	}
	return false, fmt.Sprintf("obtained value %s:%#v not supported", p0value.Kind(), params[0])
}
