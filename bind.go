// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"

	"github.com/kr/pretty"
)

// Binding provides a way to bind arguments to a checker. You can use it as a
// checker or as a matcher.
type Binding interface {
	Checker
	// Matches returns true if x passes the checker.
	Matches(x any) bool
	// String returns the name of the sub-checker and the bound arguments.
	String() string
}

// Bind takes a checker and binds arguments to the checker. It returns a Binding
// which can be used as a checker or matcher.
func Bind(check Checker, args ...any) Binding {
	if binding, ok := check.(Binding); ok && len(args) == 0 {
		return binding
	}

	info := check.Info()
	if len(info.Params) == 0 {
		panic("cannot bind checker without input parameter")
	} else if len(info.Params) <= len(args) {
		panic(fmt.Sprintf(
			"too many parameters: %s can only take %d but got %d",
			info.Name, len(info.Params)-1, len(args),
		))
	}

	b := bind{
		checker: check,
		info: CheckerInfo{
			Name:   fmt.Sprintf("%s(#%d)", info.Name, len(args)),
			Params: info.Params[:len(info.Params)-len(args)],
		},
		args: args,
	}

	return &b
}

type bind struct {
	checker Checker
	info    CheckerInfo
	args    []any
}

func (b *bind) Info() *CheckerInfo {
	return &b.info
}

func (b *bind) Check(params []any, names []string) (result bool, error string) {
	final := make([]any, 0, len(params)+len(b.args))
	final = append(final, params...)
	final = append(final, b.args...)
	return b.checker.Check(final, b.checker.Info().Params)
}

func (b *bind) Matches(x any) bool {
	final := make([]any, 0, len(b.args)+1)
	final = append(final, x)
	final = append(final, b.args...)
	match, _ := b.checker.Check(final, b.Info().Params)
	return match
}

func (b *bind) String() string {
	return fmt.Sprintf("%s(%s)", b.checker.Info().Name, pretty.Sprint(b.args...))
}
