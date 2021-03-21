// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"go/ast"
	"sort"
)

type analyzer struct {
	fields []*ast.Field
	objs   map[string]map[string]*ast.FuncDecl
}

func newAnalyzer() *analyzer {
	return &analyzer{}
}

func (a *analyzer) findFields(name string, intfs map[string]*ast.InterfaceType) []*ast.Field {
	intf, ok := intfs[name]
	if !ok {
		return nil // interface not defined in package or defined in generated file (excluded)
	}

	fields := make([]*ast.Field, 0, 10)

	for _, field := range intf.Methods.List {
		switch t := field.Type.(type) {
		case *ast.Ident:
			fields = append(fields, a.findFields(t.Name, intfs)...)
		case *ast.FuncType:
			fields = append(fields, field)
		}
	}
	sort.SliceStable(fields, func(i, j int) bool { return fields[i].Names[0].Name < fields[j].Names[0].Name })
	return fields
}

func (a *analyzer) findIntfs(node ast.Node) map[string]*ast.InterfaceType {
	intfs := make(map[string]*ast.InterfaceType, 25)

	ast.Inspect(node, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.GenDecl: // inspect generic declarations only
			ast.Inspect(node, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.TypeSpec:
					if intf, ok := node.Type.(*ast.InterfaceType); ok {
						intfs[node.Name.Name] = intf
						return false
					}
				}
				return true
			})
			return false
		}
		return true
	})
	return intfs
}

func (a *analyzer) objName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return a.objName(expr.X)
	}
	panic("cannot detect object name")
}

func (a *analyzer) findObjs(node ast.Node) map[string]map[string]*ast.FuncDecl {
	objs := make(map[string]map[string]*ast.FuncDecl, 25)

	ast.Inspect(node, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.FuncDecl:
			if node.Recv == nil || len(node.Recv.List) != 1 {
				return true
			}
			field := node.Recv.List[0]
			if len(field.Names) != 1 {
				return true
			}
			if field.Names[0].Obj == nil {
				return true
			}
			objName := a.objName(field.Type)
			mths, ok := objs[objName]
			if ok {
				mths[node.Name.Name] = node
			} else {
				objs[objName] = map[string]*ast.FuncDecl{node.Name.Name: node}
			}
			return false
		}
		return true
	})
	return objs
}

func (a *analyzer) analyze(node ast.Node) {
	intfs := a.findIntfs(node)
	a.fields = a.findFields("Converter", intfs)
	a.objs = a.findObjs(node)
}
