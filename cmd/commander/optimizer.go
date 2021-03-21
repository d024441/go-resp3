// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
)

/*
optimizer analyses generated source code and optimize it to keep code generation simple.

combine consecutive adds:
	r.add(elem1)
	r.add(elem2)
	-->
	r.add(elem1, elem2)

*/

type optimizer struct {
}

func newOptimizer() *optimizer {
	return &optimizer{}
}

func (o *optimizer) isEncode(assignStmt *ast.AssignStmt) bool {
	if len(assignStmt.Rhs) != 1 {
		return false
	}
	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != "append" {
		return false
	}
	return true
}

func (o *optimizer) mergeEncode(assignStmt1, assignStmt2 *ast.AssignStmt) {
	callExpr1 := assignStmt1.Rhs[0].(*ast.CallExpr)
	callExpr2 := assignStmt2.Rhs[0].(*ast.CallExpr)

	for i := 1; i < len(callExpr2.Args); i++ { // start with secaond argument
		callExpr1.Args = append(callExpr1.Args, callExpr2.Args[i])
	}
}

func (o *optimizer) optimize(src []byte) ([]byte, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return src, err
	}

	ast.Inspect(f, func(node ast.Node) bool {

		var stmts *[]ast.Stmt

		switch node := node.(type) {

		case *ast.BlockStmt:
			stmts = &node.List
			node.Lbrace = token.NoPos
			node.Rbrace = token.NoPos
		case *ast.CaseClause:
			stmts = &node.Body
			node.Case = token.NoPos

		case *ast.BasicLit:
			node.ValuePos = token.NoPos
		case *ast.Ident:
			node.NamePos = token.NoPos
		case *ast.IfStmt:
			node.If = token.NoPos
		case *ast.RangeStmt:
			node.For = token.NoPos
		}

		if stmts != nil && *stmts != nil {

			var lastEncode *ast.AssignStmt
			j := 0

			for _, stmt := range *stmts {
				if assignStmt, ok := stmt.(*ast.AssignStmt); ok && o.isEncode(assignStmt) {
					if lastEncode == nil {
						lastEncode = assignStmt
						(*stmts)[j] = assignStmt
						j++
					} else {
						o.mergeEncode(lastEncode, assignStmt)
					}
				} else {
					lastEncode = nil
					(*stmts)[j] = stmt
					j++
				}
			}
			*stmts = (*stmts)[:j]
		}
		return true
	})

	var b bytes.Buffer
	if err := (&printer.Config{Tabwidth: 8}).Fprint(&b, fset, f); err != nil {
		return src, err
	}

	src, err = format.Source(b.Bytes())
	if err != nil {
		return b.Bytes(), err
	}
	return src, nil
}
