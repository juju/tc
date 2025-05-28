// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"crypto/rand"
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
	"sync"
)

// MultiChecker is a deep checker that by default matches for equality.
// But checks can be overriden based on path (either explicit match or regexp)
type MultiChecker struct {
	*CheckerInfo
	matchChecks []matchCheck
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

// NewMultiChecker creates a MultiChecker which is a deep checker that by default matches for equality.
// But checks can be overriden based on path (either explicit match or regexp)
func NewMultiChecker() *MultiChecker {
	return &MultiChecker{
		CheckerInfo: &CheckerInfo{Name: "MultiChecker", Params: []string{"obtained", "expected"}},
	}
}

// AddExpr exception which matches path with go expression. Use _ for wildcard.
// The top level or root value must be a _ when using expression.
func (checker *MultiChecker) AddExpr(expr string, c Checker, args ...any) *MultiChecker {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	root := findRoot(astExpr, "_")
	if root == nil {
		panic("cannot find root ident _")
	}
	root.Name = topLevel

	astExpr = simplify(astExpr)

	checker.matchChecks = append(checker.matchChecks, &astCheck{
		multiCheck: multiCheck{
			Checker: c,
			args:    args,
		},
		astExpr: astExpr,
	})
	return checker
}

// topLevel is a substitute for the top level or root object.
// We use an unlikely value to provide backwards compatability with previous deep equals
// behaviour. It is stripped out before any errors are printed.
var topLevel = "_" + rand.Text()

// Check for go check Checker interface.
func (checker *MultiChecker) Check(params []any, names []string) (result bool, errStr string) {
	customCheckFunc := func(path string, a1 any, a2 any) (useDefault bool, equal bool, err error) {
		var mc checkerWithArgs
		for _, v := range checker.matchChecks {
			if v.MatchString(path) {
				mc = v
				break
			}
		}
		if mc == nil {
			return true, false, nil
		}

		params := append([]any{a1}, mc.Args()...)
		info := mc.Info()
		if len(params) < len(info.Params) {
			return false, false, fmt.Errorf("Wrong number of parameters for %s: want %d, got %d", info.Name, len(info.Params), len(params)+1)
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
			return false, true, nil
		}
		path = strings.Replace(path, topLevel, "", 1)
		if path == "" {
			path = "top level"
		}
		return false, false, fmt.Errorf("mismatch at %s: %s", path, errStr)
	}
	if ok, err := DeepEqualWithCustomCheck(params[0], params[1], customCheckFunc); !ok {
		return false, err.Error()
	}
	return true, ""
}

// ExpectedValue if passed to MultiChecker.AddExpr, will be substituded with the expected value.
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
