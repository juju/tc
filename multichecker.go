// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"reflect"
	"sync"

	"github.com/kr/pretty"
)

// MultiChecker is a deep checker that by default matches for equality.
// But checks can be overriden based on path (either explicit match or regexp)
type MultiChecker struct {
	*CheckerInfo
	matchChecks       []matchCheck
	lengthMatchChecks []matchCheck
	equals            checkerWithArgs
}

type checkerWithArgs interface {
	Checker
	Args() []any
}

type matchCheck interface {
	checkerWithArgs
	MatchString(string) bool
	WantTopLevel() bool
}

type multiCheck struct {
	Checker
	args []any
}

func (m *multiCheck) Args() []any {
	return m.args
}

type astCheck struct {
	multiCheck
	astExpr ast.Expr
}

func (a *astCheck) WantTopLevel() bool {
	return true
}

// NewMultiChecker creates a MultiChecker which is a deep checker that by
// default matches for equality. But checks can be overriden based on path
// (either explicit match or regexp)
func NewMultiChecker() *MultiChecker {
	return &MultiChecker{
		CheckerInfo: &CheckerInfo{Name: "MultiChecker", Params: []string{"obtained", "expected"}},
	}
}

// SetDefault changes the default equality checking to another checker.
// Using [tc.Ignore] or [tc.Equals] (with [tc.ExpectedValue]) is currently the
// only reasonable checkers.
func NewMultiCheckerWithDefault(c Checker, args ...any) *MultiChecker {
	mc := NewMultiChecker()
	mc.equals = &multiCheck{
		Checker: c,
		args:    args,
	}
	return mc
}

// AddExpr exception which matches path with go expression. Use _ for wildcard.
// The top level or root value must be a _ when using expression.
// Use `len(_.x.y.z)` to override length checking for the path.
func (checker *MultiChecker) AddExpr(
	expr string, c Checker, args ...any,
) *MultiChecker {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	isLenChecker := false
	if callExpr, ok := astExpr.(*ast.CallExpr); ok {
		if lhs, ok := callExpr.Fun.(*ast.Ident); !ok || lhs.Name != "len" {
			panic(fmt.Errorf(
				"call expression only supports len: got %s",
				pretty.Sprint(callExpr.Fun),
			))
		}
		if len(callExpr.Args) != 1 {
			panic("len call expression expected 1 argument")
		}
		astExpr = callExpr.Args[0]
		isLenChecker = true
	}

	root := findRoot(astExpr, "_")
	if root == nil {
		panic("cannot find root ident _")
	}
	root.Name = topLevel

	astExpr = simplify(astExpr)

	astChecker := &astCheck{
		multiCheck: multiCheck{
			Checker: c,
			args:    args,
		},
		astExpr: astExpr,
	}

	if isLenChecker {
		checker.lengthMatchChecks = append(checker.lengthMatchChecks, astChecker)
	} else {
		checker.matchChecks = append(checker.matchChecks, astChecker)
	}
	return checker
}

// topLevel is a substitute for the top level or root object.
// We use an unlikely value to provide backwards compatability with previous deep equals
// behaviour. It is stripped out before any errors are printed.
var topLevel = "_" + rand.Text()

// Check for go check Checker interface.
func (checker *MultiChecker) Check(
	params []any, names []string,
) (result bool, errStr string) {
	v1 := reflect.ValueOf(params[0])
	v2 := reflect.ValueOf(params[1])
	result, err := deepValueEqual(topLevel, v1, v2, make(map[visit]bool), 0,
		checker.customEquals,
		checker.customLength,
		checker.customCheck)
	if err != nil {
		return result, err.Error()
	}
	return result, ""
}

func (checker *MultiChecker) customCheck(
	path string, a1 any, a2 any,
) (useDefault bool, equal bool, err error) {
	var checkers []checkerWithArgs
	for _, v := range checker.matchChecks {
		if v.MatchString(path) {
			checkers = append(checkers, v)
		}
	}
	if len(checkers) == 0 {
		return true, false, nil
	}

	for _, mc := range checkers {
		params := append([]any{a1}, mc.Args()...)
		info := mc.Info()
		if len(params) < len(info.Params) {
			return false, false, fmt.Errorf(
				"Wrong number of parameters for %s: want %d, got %d",
				info.Name, len(info.Params), len(params),
			)
		}
		// Copy since it may be mutated by Check.
		names := append([]string{}, info.Params...)

		// Trim to the expected params len.
		params = params[:len(info.Params)]

		// Perform substitution
		for i, v := range params {
			if v == ExpectedValue {
				params[i] = a2
			}
		}

		result, errStr := mc.Check(params, names)
		if result {
			continue
		}

		var err error
		if errStr != "" {
			err = errors.New(errStr)
		}
		return false, false, err
	}

	return false, true, nil
}

func (checker *MultiChecker) customEquals(a1 any, a2 any) bool {
	if checker.equals == nil {
		return a1 == a2
	}

	params := append([]any{a1}, checker.equals.Args()...)
	info := checker.equals.Info()
	if len(params) < len(info.Params) {
		panic(fmt.Errorf(
			"Wrong number of parameters for %s: want %d, got %d",
			info.Name, len(info.Params), len(params),
		))
	}
	// Copy since it may be mutated by Check.
	names := append([]string{}, info.Params...)

	// Trim to the expected params len.
	params = params[:len(info.Params)]

	// Perform substitution
	for i, v := range params {
		if v == ExpectedValue {
			params[i] = a2
		}
	}

	result, _ := checker.equals.Check(params, names)
	return result
}

func (checker *MultiChecker) customLength(path string, a1 int, a2 int) bool {
	var checkers []checkerWithArgs
	for _, v := range checker.lengthMatchChecks {
		if v.MatchString(path) {
			checkers = append(checkers, v)
		}
	}
	if len(checkers) == 0 {
		return a1 == a2
	}

	for _, mc := range checkers {
		params := append([]any{a1}, mc.Args()...)
		info := mc.Info()
		if len(params) < len(info.Params) {
			panic(fmt.Errorf(
				"Wrong number of parameters for %s: want %d, got %d",
				info.Name, len(info.Params), len(params)+1,
			))
		}
		// Copy since it may be mutated by Check.
		names := append([]string{}, info.Params...)

		// Trim to the expected params len.
		params = params[:len(info.Params)]

		// Perform substitution
		for i, v := range params {
			if v == ExpectedValue {
				params[i] = a2
			}
		}

		result, _ := mc.Check(params, names)
		if result {
			continue
		}

		return false
	}

	return true
}

// ExpectedValue if passed to MultiChecker.AddExpr, will be substituded with the
// expected value.
var ExpectedValue = &struct{}{}

var (
	astCache     = make(map[string]ast.Expr)
	astCacheLock = sync.Mutex{}
)

func (a *astCheck) MatchString(expr string) bool {
	astCacheLock.Lock()
	astExpr, ok := astCache[expr]
	astCacheLock.Unlock()
	if !ok {
		var err error
		astExpr, err = parser.ParseExpr(expr)
		if err != nil {
			panic(err)
		}
		astExpr = simplify(astExpr)
		astCacheLock.Lock()
		astCache[expr] = astExpr
		astCacheLock.Unlock()
	}

	if matchAstExpr(a.astExpr, astExpr) {
		return true
	} else {
		return false
	}
}

func simplify(x ast.Expr) ast.Expr {
	switch expr := x.(type) {
	case *ast.IndexExpr:
		copyExpr := *expr
		copyExpr.X = simplify(expr.X)
		copyExpr.Index = simplify(expr.Index)
		return &copyExpr
	case *ast.ParenExpr:
		return simplify(expr.X)
	case *ast.StarExpr:
		return simplify(expr.X)
	case *ast.SelectorExpr:
		exprCopy := *expr
		exprCopy.X = simplify(expr.X)
		return &exprCopy
	case *ast.Ident:
		return expr
	case *ast.BasicLit:
		return expr
	default:
		panic(fmt.Sprintf("unknown type %#v", expr))
	}
}

func findRoot(x ast.Expr, name string) *ast.Ident {
	switch expr := x.(type) {
	case *ast.IndexExpr:
		return findRoot(expr.X, name)
	case *ast.ParenExpr:
		return findRoot(expr.X, name)
	case *ast.StarExpr:
		return findRoot(expr.X, name)
	case *ast.SelectorExpr:
		return findRoot(expr.X, name)
	case *ast.Ident:
		if expr.Name == name {
			return expr
		}
	case *ast.BasicLit:
	default:
		panic(fmt.Sprintf("unknown type %#v", expr))
	}
	return nil
}

func matchAstExpr(expected, obtained ast.Expr) bool {
	switch expr := expected.(type) {
	case *ast.IndexExpr:
		x, ok := obtained.(*ast.IndexExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
		if !matchAstExpr(expr.Index, x.Index) {
			return false
		}
	case *ast.ParenExpr:
		x, ok := obtained.(*ast.ParenExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
	case *ast.StarExpr:
		x, ok := obtained.(*ast.StarExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
	case *ast.SelectorExpr:
		x, ok := obtained.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if !matchAstExpr(expr.X, x.X) {
			return false
		}
		if !matchAstExpr(expr.Sel, x.Sel) {
			return false
		}
	case *ast.Ident:
		if expr.Name == "_" {
			// Wildcard
			return true
		}
		x, ok := obtained.(*ast.Ident)
		if !ok {
			return false
		}
		if expr.Name != x.Name {
			return false
		}
	case *ast.BasicLit:
		x, ok := obtained.(*ast.BasicLit)
		if !ok {
			return false
		}
		if expr.Kind != x.Kind {
			return false
		}
		if expr.Value != x.Value {
			return false
		}
	default:
		panic(fmt.Sprintf("unknown type %#v", expected))
	}
	return true
}
