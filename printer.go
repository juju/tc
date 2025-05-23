// Gocheck - A rich testing framework for Go
//
// Copyright (c) 2010-2013 Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package tc

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func indent(s, with string) (r string) {
	eol := true
	for i := 0; i != len(s); i++ {
		c := s[i]
		switch {
		case eol && c == '\n' || c == '\r':
		case c == '\n' || c == '\r':
			eol = true
		case eol:
			eol = false
			s = s[:i] + with + s[i:]
			i += len(with)
		}
	}
	return s
}

func printLine(filename string, line int) (string, error) {
	fset := token.NewFileSet()
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	fnode, err := parser.ParseFile(fset, filename, file, parser.ParseComments)
	if err != nil {
		return "", err
	}
	config := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	lp := &linePrinter{fset: fset, fnode: fnode, line: line, config: config}
	ast.Walk(lp, fnode)
	result := lp.output.Bytes()
	// Comments leave \n at the end.
	n := len(result)
	for n > 0 && result[n-1] == '\n' {
		n--
	}
	return string(result[:n]), nil
}

type linePrinter struct {
	config *printer.Config
	fset   *token.FileSet
	fnode  *ast.File
	line   int
	output bytes.Buffer
	stmt   ast.Stmt
}

func (lp *linePrinter) emit() bool {
	if lp.stmt != nil {
		lp.trim(lp.stmt)
		lp.printWithComments(lp.stmt)
		lp.stmt = nil
		return true
	}
	return false
}

func (lp *linePrinter) printWithComments(n ast.Node) {
	nfirst := lp.fset.Position(n.Pos()).Line
	nlast := lp.fset.Position(n.End()).Line
	for _, g := range lp.fnode.Comments {
		cfirst := lp.fset.Position(g.Pos()).Line
		clast := lp.fset.Position(g.End()).Line
		if clast == nfirst-1 && lp.fset.Position(n.Pos()).Column == lp.fset.Position(g.Pos()).Column {
			for _, c := range g.List {
				lp.output.WriteString(c.Text)
				lp.output.WriteByte('\n')
			}
		}
		if cfirst >= nfirst && cfirst <= nlast && n.End() <= g.List[0].Slash {
			// The printer will not include the comment if it starts past
			// the node itself. Trick it into printing by overlapping the
			// slash with the end of the statement.
			g.List[0].Slash = n.End() - 1
		}
	}
	node := &printer.CommentedNode{
		Node:     n,
		Comments: lp.fnode.Comments,
	}
	lp.config.Fprint(&lp.output, lp.fset, node)
}

func (lp *linePrinter) Visit(n ast.Node) (w ast.Visitor) {
	if n == nil {
		if lp.output.Len() == 0 {
			lp.emit()
		}
		return nil
	}
	first := lp.fset.Position(n.Pos()).Line
	last := lp.fset.Position(n.End()).Line
	if first <= lp.line && last >= lp.line {
		// Print the innermost statement containing the line.
		if stmt, ok := n.(ast.Stmt); ok {
			if _, ok := n.(*ast.BlockStmt); !ok {
				lp.stmt = stmt
			}
		}
		if first == lp.line && lp.emit() {
			return nil
		}
		return lp
	}
	return nil
}

func (lp *linePrinter) trim(n ast.Node) bool {
	stmt, ok := n.(ast.Stmt)
	if !ok {
		return true
	}
	line := lp.fset.Position(n.Pos()).Line
	if line != lp.line {
		return false
	}
	switch stmt := stmt.(type) {
	case *ast.IfStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.SwitchStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.TypeSwitchStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.CaseClause:
		stmt.Body = lp.trimList(stmt.Body)
	case *ast.CommClause:
		stmt.Body = lp.trimList(stmt.Body)
	case *ast.BlockStmt:
		stmt.List = lp.trimList(stmt.List)
	}
	return true
}

func (lp *linePrinter) trimBlock(stmt *ast.BlockStmt) *ast.BlockStmt {
	if !lp.trim(stmt) {
		return lp.emptyBlock(stmt)
	}
	stmt.Rbrace = stmt.Lbrace
	return stmt
}

func (lp *linePrinter) trimList(stmts []ast.Stmt) []ast.Stmt {
	for i := 0; i != len(stmts); i++ {
		if !lp.trim(stmts[i]) {
			stmts[i] = lp.emptyStmt(stmts[i])
			break
		}
	}
	return stmts
}

func (lp *linePrinter) emptyStmt(n ast.Node) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.Ellipsis{
			Ellipsis: n.Pos(),
		},
	}
}

func (lp *linePrinter) emptyBlock(n ast.Node) *ast.BlockStmt {
	p := n.Pos()
	return &ast.BlockStmt{
		Lbrace: p,
		List:   []ast.Stmt{lp.emptyStmt(n)},
		Rbrace: p,
	}
}
