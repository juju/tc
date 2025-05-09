// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"reflect"
	"slices"
)

type orderedChecker[T ~[]E, E any] struct {
	*CheckerInfo
	matcher Checker
	right   bool
}

// OrderedLeft checks the left/obtained value is a slice
// and all of its values appear in the right/expected slice
// in the same order. A matcher is passed to match values
// between the slices.
func OrderedLeft[T ~[]E, E any](matcher Checker) Checker {
	return &orderedChecker[T, E]{
		CheckerInfo: &CheckerInfo{
			Name:   fmt.Sprintf("OrderedLeft[%s]", reflect.TypeFor[T]().Name()),
			Params: []string{"obtained", "expected"},
		},
		matcher: matcher,
		right:   false,
	}
}

// OrderedRight checks the left/obtained value is a slice
// and all of its values appear in the right/expected slice
// in the same order. A matcher is passed to match values
// between the slices.
func OrderedRight[T ~[]E, E any](matcher Checker) Checker {
	return &orderedChecker[T, E]{
		CheckerInfo: &CheckerInfo{
			Name:   fmt.Sprintf("OrderedRight[%s]", reflect.TypeFor[T]().Name()),
			Params: []string{"obtained", "expected"},
		},
		matcher: matcher,
		right:   true,
	}
}

func (o *orderedChecker[T, E]) Info() *CheckerInfo {
	return o.CheckerInfo
}

func (o *orderedChecker[T, E]) Check(params []any, names []string) (result bool, error string) {
	if len(params) != 2 {
		return false, o.Name + " expects two typed slice arguments"
	}
	obtained, ok := params[0].(T)
	if !ok {
		return false, fmt.Sprintf("%s expects left type %s, got %s",
			o.Name,
			reflect.TypeFor[T]().Name(),
			reflect.TypeOf(params[0]).Name())
	}
	expected := params[1].(T)
	if !ok {
		return false, fmt.Sprintf("%s expects right type %s, got %s",
			o.Name,
			reflect.TypeFor[T]().Name(),
			reflect.TypeOf(params[1]).Name())
	}

	var want T
	var have T
	if o.right {
		want = slices.Clone(expected)
		have = obtained
	} else {
		want = slices.Clone(obtained)
		have = expected
	}
	for _, v := range have {
		if len(want) == 0 {
			break
		}
		var values []any
		if o.right {
			values = []any{v, want[0]}
		} else {
			values = []any{want[0], v}
		}
		if matches, _ := o.matcher.Check(values, o.matcher.Info().Params); matches {
			want = slices.Delete(want, 0, 1)
		}
	}

	if len(want) != 0 {
		return false, fmt.Sprintf("%d unmatched elements", len(want))
	}

	return true, ""
}
