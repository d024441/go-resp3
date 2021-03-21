// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/stfnmllr/go-resp3/cmd/commander/internal/ast"
	"github.com/stfnmllr/go-resp3/cmd/commander/internal/stringutils"
)

var cmdNameMap = map[string]string{
	"TTL":      "TTL",
	"PTTL":     "PTTL",
	"HSETNX":   "HsetNx",
	"RENAMENX": "RenameNx",
	"MSETNX":   "MsetNx",
}

const (
	argNameIDOrDollar = "id-or-$"
)

var argNameMap = map[string]string{
	argNameIDOrDollar: "id",
}

const (
	typeString    = "string"
	typeKey       = "key"
	typeType      = "type"
	typeInteger   = "integer"
	typeDouble    = "double"
	typePattern   = "pattern"
	typeChannel   = "channel"
	typePosixTime = "posix time"
)

var argTypeMap = map[string]ast.TypeNode{
	typeString:    ast.StringType,
	typeKey:       ast.KeyType,
	typeType:      ast.StringType,
	typeInteger:   ast.IntegerType,
	typeDouble:    ast.FloatType,
	typePattern:   ast.StringType,
	typeChannel:   ast.StringType,
	typePosixTime: ast.IntegerType,
}

// replace type of map[groupName][typeName]type
var specialTypeMap = map[string]map[string]ast.TypeNode{
	"geo":         {"member": ast.InterfaceType, "member1": ast.InterfaceType, "member2": ast.InterfaceType},
	"hash":        {"field": ast.InterfaceType, "value": ast.InterfaceType},
	"hyperloglog": {"element": ast.InterfaceType},
	"list":        {"pivot": ast.InterfaceType, "element": ast.InterfaceType},
	"scripting":   {"arg": ast.InterfaceType},
	"set":         {"member": ast.InterfaceType},
	"sorted_set":  {"member": ast.InterfaceType},
	"stream":      {"field": ast.InterfaceType, "value": ast.InterfaceType},
	"string":      {"value": ast.InterfaceType},
}

type converter struct {
	s *ast.Scope
}

func newConverter(s *ast.Scope) *converter {
	return &converter{s: s}
}

func (c *converter) convert(commands map[string]*command) {
	for cmd, command := range commands {
		c.adapt(cmd, command)
		c.convertCommand(cmd, command)
	}
}

func (c *converter) adapt(cmd string, command *command) {
	switch cmd {
	case "MIGRATE":
		for _, a := range command.Arguments {
			if a.Name.id == "key" && len(a.Enum) != 0 { // correct enum to key type
				a.Type.id = "key"
				a.Enum = nil
			}
		}
	}
}

func (c *converter) convertCommand(cmdKey string, cmd *command) {
	name := cmdName(cmdKey)
	funcDecl := ast.NewFuncDecl(name, name, strings.Split(cmdKey, " "))
	funcAttr := ast.NewFuncAttr(name, cmd.Summary, cmd.Complexity, cmd.Since, cmd.Group)

	c.s.InsertDecl(funcAttr)
	if !c.s.InsertDecl(funcDecl) {
		return // declaration was provided by patch file
	}

	for _, a := range cmd.Arguments {

		switch a.Kind() {

		case akBasic:
			field := &ast.Field{Name: convertPrmName(a), Cmd: a.Command, Type: convertType(cmd.Group, a)}
			funcDecl.List = append(funcDecl.List, field)

		case akStruct:
			structDecl := ast.NewStructDecl(convertStructTypeName(a))
			for i := 0; i < len(a.Name.ids); i++ {
				field := &ast.Field{Name: convertElemNameAt(i, a), Type: convertTypeAt(cmd.Group, i, a)}
				structDecl.List = append(structDecl.List, field)
			}
			decl := c.s.Lookup(structDecl.Name)
			if decl == nil {
				c.s.InsertDecl(structDecl)
			} else {
				d, ok := decl.(*ast.StructDecl)
				if !ok {
					log.Fatalf("function %s invalid object type %T - expected %T", funcDecl.Name, structDecl, decl)
				}
				if !reflect.DeepEqual(d, structDecl) {
					log.Fatalf("function %s struct type mismatch %#v expected: %#v", funcDecl.Name, structDecl, d)
				}
			}
			field := &ast.Field{Name: convertPrmStructName(a), Cmd: a.Command, Type: convertStructType(a)}
			funcDecl.List = append(funcDecl.List, field)

		case akEnum:
			enumDecl := ast.NewEnumDecl(convertEnumTypeName(a), a.Enum)
			decl := c.s.Lookup(enumDecl.Name)
			if decl == nil {
				c.s.InsertDecl(enumDecl)
			} else {
				d, ok := decl.(*ast.EnumDecl)
				if !ok {
					log.Fatalf("invalid object type %T - expected %T", enumDecl, decl)
				}
				if !reflect.DeepEqual(decl, enumDecl) {
					log.Fatalf("enum type mismatch %s expected: %s ", enumDecl, d)
				}
			}
			field := &ast.Field{Name: convertPrmName(a), Cmd: a.Command, Type: convertEnumType(a)}
			funcDecl.List = append(funcDecl.List, field)

		case akEnumBool:
			field := &ast.Field{Name: convertPrmName(a), Type: convertEnumBoolType(a)}
			funcDecl.List = append(funcDecl.List, field)

		case akEnumConst:
			field := &ast.Field{Name: convertPrmName(a), Cmd: a.Enum[0]}
			funcDecl.List = append(funcDecl.List, field)

		}

	}
}

func cmdName(name string) string {
	if r, ok := cmdNameMap[name]; ok {
		return r
	}
	return stringutils.PascalCase(name)
}

func normName(name string) string {
	if r, ok := argNameMap[name]; ok {
		return r
	}
	return name
}

func normNames(names []string) []string {
	r := make([]string, len(names))
	for i, name := range names {
		r[i] = normName(name)
	}
	return r
}

func normType(group, name, typ string) ast.TypeNode {
	r, ok := argTypeMap[typ]
	if !ok {
		log.Fatalf("cannot convert name %s type %s", name, typ)
	}
	// special type conversions
	if r == ast.StringType {
		if s, ok := specialTypeMap[group][name]; ok {
			return s
		}
	}
	return r
}

func convertPrmName(a *argument) string {
	var r string
	switch {
	// if two value enum -> take first enum value as name
	case a.Kind() == akEnumBool && len(a.Enum) == 2: // if two value enum -> take first enum value as name
		r = stringutils.CamelCase(normName(a.Enum[0]))
	case a.Command != "":
		r = stringutils.CamelCase(a.Command)
	default:
		r = stringutils.CamelCase(normName(a.Name.id))
	}
	if r == "type" { // edge case - cannot use "type" as parameter name
		return "typ"
	}
	return r
}

func convertType(group string, a *argument) ast.TypeNode {
	typ := normType(group, a.Name.id, a.Type.id)
	switch {
	case a.Multiple || a.Variadic:
		typ = &ast.SliceType{AllowNil: a.Optional, Variadic: a.Variadic, Node: typ}
	case a.Optional:
		typ = &ast.PointerType{Node: typ}
	}
	return typ
}

func convertElemNameAt(i int, a *argument) string {
	return stringutils.PascalCase(normName(a.Name.ids[i]))
}

func convertTypeAt(group string, i int, a *argument) ast.TypeNode {
	return normType(group, a.Name.ids[i], a.Type.ids[i])
}

func convertPrmStructName(a *argument) string {
	if a.Command != "" {
		return stringutils.CamelCase(a.Command)
	}
	return stringutils.CamelCase(strings.Join(normNames(a.Name.ids), " "))
}

func convertStructTypeName(a *argument) string {
	return stringutils.PascalCase(strings.Join(normNames(a.Name.ids), " "))
}

func convertStructType(a *argument) ast.TypeNode {
	var typ ast.TypeNode

	typ = &ast.DataType{Name: convertStructTypeName(a)}
	switch {
	case a.Multiple || a.Variadic:
		typ = &ast.SliceType{Node: typ}
	case a.Optional:
		typ = &ast.PointerType{Node: typ}
	}
	return typ
}

func convertEnumTypeName(a *argument) string {
	if a.Name.id != "" {
		return stringutils.PascalCase(normName(a.Name.id))
	}
	return stringutils.PascalCase(a.Command)
}

func convertEnumType(a *argument) ast.TypeNode {
	var typ ast.TypeNode

	typ = &ast.DataType{Name: convertEnumTypeName(a)}
	switch {
	case a.Multiple || a.Variadic:
		typ = &ast.SliceType{Node: typ}
	case a.Optional:
		typ = &ast.PointerType{Node: typ}
	}
	return typ
}

func convertEnumBoolType(a *argument) ast.TypeNode {
	var typ ast.TypeNode

	typ = &ast.EnumBoolType{Values: a.Enum}
	if len(a.Enum) != 1 {
		switch {
		case a.Multiple || a.Variadic:
			typ = &ast.SliceType{Node: typ}
		case a.Optional:
			typ = &ast.PointerType{Node: typ}
		}
	}
	return typ
}
