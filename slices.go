// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/kr/pretty"
)

type orderedChecker[T ~[]E, E any] struct {
	*CheckerInfo
	matcher Checker
	right   bool
	full    bool
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

// OrderedMatch checks the obtained slice contains the
// same values as the expected, in the same order.
func OrderedMatch[T ~[]E, E any](matcher Checker) Checker {
	return &orderedChecker[T, E]{
		CheckerInfo: &CheckerInfo{
			Name:   fmt.Sprintf("OrderedMatch[%s]", reflect.TypeFor[T]().Name()),
			Params: []string{"obtained", "expected"},
		},
		matcher: matcher,
		full:    true,
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
	for i, v := range have {
		matched := false
		if len(want) > 0 {
			var values []any
			if o.right {
				values = []any{v, want[0]}
			} else {
				values = []any{want[0], v}
			}
			if matches, _ := o.matcher.Check(values, o.matcher.Info().Params); matches {
				want = slices.Delete(want, 0, 1)
				matched = true
			}
		}
		if !matched && o.full {
			if o.right {
				return false, fmt.Sprintf("unexpected element: %s", pretty.Sprint(have[i]))
			} else {
				return false, fmt.Sprintf("expected elements missing: %s", pretty.Sprint(have[i:]))
			}
		}
	}

	if len(want) != 0 {
		if o.right {
			return false, fmt.Sprintf("expected elements missing: %s", pretty.Sprint(want))
		} else {
			return false, fmt.Sprintf("unexpected element: %s", pretty.Sprint(want[0]))
		}
	}

	return true, ""
}

type unorderedChecker[T ~[]E, E any] struct {
	*CheckerInfo
	matcher Checker
}

// UnorderedMatch checks the obtained slice contains the
// same values as the expected, but in any order.
func UnorderedMatch[T ~[]E, E any](matcher Checker) Checker {
	return &unorderedChecker[T, E]{
		CheckerInfo: &CheckerInfo{
			Name:   fmt.Sprintf("UnorderedMatch[%s]", reflect.TypeFor[T]().Name()),
			Params: []string{"obtained", "expected"},
		},
		matcher: matcher,
	}
}

func (o *unorderedChecker[T, E]) Info() *CheckerInfo {
	return o.CheckerInfo
}

func (o *unorderedChecker[T, E]) Check(params []any, names []string) (result bool, error string) {
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

	obtained = slices.Clone(obtained)
	for _, right := range expected {
		matched := false
		for k, left := range obtained {
			if matches, _ := o.matcher.Check([]any{left, right}, o.matcher.Info().Params); matches {
				obtained = slices.Delete(obtained, k, k+1)
				matched = true
				break
			}
		}
		if !matched {
			return false, fmt.Sprintf("expected element missing: %s", pretty.Sprint(right))
		}
	}

	if len(obtained) != 0 {
		return false, fmt.Sprintf("%d unmatched elements: %s", len(obtained), pretty.Sprint(obtained))
	}

	return true, ""
}
