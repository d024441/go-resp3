// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package ast

import (
	"encoding/json"
	"reflect"
)

//go:generate stringer -type=nodeKind

type nodeKind int

const (
	unknown nodeKind = iota
	funcDecl
	funcAttr
	funcConfig
	structDecl
	enumDecl
	alternative
	field
	baseType
	dataType
	enumBoolType
	pointerType
	sliceType
	ellipsisType
	cbType
	maxNodeKind
)

var nodeKindString = map[string]nodeKind{}

func init() {
	for i := 0; i < int(maxNodeKind); i++ {
		nodeKindString[nodeKind(i).String()] = nodeKind(i)
	}
}

type jsonNode interface {
	MarshalJSON() ([]byte, error)
	kind() nodeKind
}

// check if Node interface is implemented
var _ jsonNode = (*FuncDecl)(nil)
var _ jsonNode = (*FuncAttr)(nil)
var _ jsonNode = (*StructDecl)(nil)
var _ jsonNode = (*EnumDecl)(nil)
var _ jsonNode = (*Alternative)(nil)
var _ jsonNode = (*Field)(nil)
var _ jsonNode = (*BaseType)(nil)
var _ jsonNode = (*DataType)(nil)
var _ jsonNode = (*EnumBoolType)(nil)
var _ jsonNode = (*PointerType)(nil)
var _ jsonNode = (*SliceType)(nil)
var _ jsonNode = (*EllipsisType)(nil)

// MarshalJSON implements the Marshaler interface.
func MarshalJSON(n jsonNode) ([]byte, error) {
	typ := reflect.TypeOf(n)
	aliasType := aliasMap[typ]

	alias := reflect.ValueOf(n).Convert(aliasType)

	node := reflect.New(jsonMap[typ]).Elem()

	node.Field(0).SetString(n.kind().String()) // set type field
	node.Field(1).Set(alias)                   // set alias

	return json.Marshal(node.Interface())
}

// Alias (need alias type - original type would lead to endles recursion in MarshalJSON

// FuncDeclAlias is the alias type for FuncDecl.
type FuncDeclAlias FuncDecl

// FuncAttrAlias is the alias type for FuncAttr.
type FuncAttrAlias FuncAttr

// FuncConfigAlias is the alias type for FuncConfig.
type FuncConfigAlias FuncConfig

// StructDeclAlias is the alias type for StructDecl.
type StructDeclAlias StructDecl

// EnumDeclAlias is the alias type for EnumDecl.
type EnumDeclAlias EnumDecl

// AlternativeAlias is the alias type for Alternative.
type AlternativeAlias Alternative

// FieldAlias is the alias type for Field.
type FieldAlias Field

// BaseTypeAlias is the alias type for BaseType.
type BaseTypeAlias BaseType

// DataTypeAlias is the alias type for DataType.
type DataTypeAlias DataType

// EnumBoolTypeAlias is the alias type for EnumBoolType.
type EnumBoolTypeAlias EnumBoolType

// PointerTypeAlias is the alias type for PointerType.
type PointerTypeAlias PointerType

// SliceTypeAlias is the alias type for SliceType.
type SliceTypeAlias SliceType

// EllipsisTypeAlias is the alias type for EllipsisType.
type EllipsisTypeAlias EllipsisType

var aliasMap = map[reflect.Type]reflect.Type{
	reflect.TypeOf((*FuncDecl)(nil)):     reflect.TypeOf((*FuncDeclAlias)(nil)),
	reflect.TypeOf((*FuncAttr)(nil)):     reflect.TypeOf((*FuncAttrAlias)(nil)),
	reflect.TypeOf((*FuncConfig)(nil)):   reflect.TypeOf((*FuncConfigAlias)(nil)),
	reflect.TypeOf((*StructDecl)(nil)):   reflect.TypeOf((*StructDeclAlias)(nil)),
	reflect.TypeOf((*EnumDecl)(nil)):     reflect.TypeOf((*EnumDeclAlias)(nil)),
	reflect.TypeOf((*Alternative)(nil)):  reflect.TypeOf((*AlternativeAlias)(nil)),
	reflect.TypeOf((*Field)(nil)):        reflect.TypeOf((*FieldAlias)(nil)),
	reflect.TypeOf((*BaseType)(nil)):     reflect.TypeOf((*BaseTypeAlias)(nil)),
	reflect.TypeOf((*DataType)(nil)):     reflect.TypeOf((*DataTypeAlias)(nil)),
	reflect.TypeOf((*EnumBoolType)(nil)): reflect.TypeOf((*EnumBoolTypeAlias)(nil)),
	reflect.TypeOf((*PointerType)(nil)):  reflect.TypeOf((*PointerTypeAlias)(nil)),
	reflect.TypeOf((*SliceType)(nil)):    reflect.TypeOf((*SliceTypeAlias)(nil)),
	reflect.TypeOf((*EllipsisType)(nil)): reflect.TypeOf((*EllipsisTypeAlias)(nil)),
}

// json structs including type information
var jsonMap = map[reflect.Type]reflect.Type{}

// create json structs, e.g.

// type jsonFuncDecl struct {
// 	Kind string `json:"_type"`
// 	*FuncDeclAlias
// }
func init() {
	for typ, alias := range aliasMap {
		jsonType := reflect.StructOf([]reflect.StructField{
			{Name: "Kind", Type: reflect.TypeOf(string("")), Tag: `json:"_type"`},
			{Name: "Dummy", Type: alias, Anonymous: true},
		})
		jsonMap[typ] = jsonType
	}
}

// Node methods

// MarshalJSON implements the Marshaler interface.
func (d *FuncDecl) MarshalJSON() ([]byte, error) { return MarshalJSON(d) }

// MarshalJSON implements the Marshaler interface.
func (d *FuncAttr) MarshalJSON() ([]byte, error) { return MarshalJSON(d) }

// MarshalJSON implements the Marshaler interface.
func (d *FuncConfig) MarshalJSON() ([]byte, error) { return MarshalJSON(d) }

// MarshalJSON implements the Marshaler interface.
func (d *StructDecl) MarshalJSON() ([]byte, error) { return MarshalJSON(d) }

// MarshalJSON implements the Marshaler interface.
func (d *EnumDecl) MarshalJSON() ([]byte, error) { return MarshalJSON(d) }

// MarshalJSON implements the Marshaler interface.
func (alt *Alternative) MarshalJSON() ([]byte, error) { return MarshalJSON(alt) }

// MarshalJSON implements the Marshaler interface.
func (f *Field) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

// MarshalJSON implements the Marshaler interface.
func (t *BaseType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// MarshalJSON implements the Marshaler interface.
func (t *DataType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// MarshalJSON implements the Marshaler interface.
func (t *EnumBoolType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// MarshalJSON implements the Marshaler interface.
func (t *PointerType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// MarshalJSON implements the Marshaler interface.
func (t *SliceType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// MarshalJSON implements the Marshaler interface.
func (t *EllipsisType) MarshalJSON() ([]byte, error) { return MarshalJSON(t) }

// Kind
func (d *FuncDecl) kind() nodeKind      { return funcDecl }
func (d *FuncAttr) kind() nodeKind      { return funcAttr }
func (d *FuncConfig) kind() nodeKind    { return funcConfig }
func (d *StructDecl) kind() nodeKind    { return structDecl }
func (d *EnumDecl) kind() nodeKind      { return enumDecl }
func (alt *Alternative) kind() nodeKind { return alternative }
func (f *Field) kind() nodeKind         { return field }
func (t *BaseType) kind() nodeKind      { return baseType }
func (t *DataType) kind() nodeKind      { return dataType }
func (t *EnumBoolType) kind() nodeKind  { return enumBoolType }
func (t *PointerType) kind() nodeKind   { return pointerType }
func (t *SliceType) kind() nodeKind     { return sliceType }
func (t *EllipsisType) kind() nodeKind  { return ellipsisType }

// Unmarshal

// NodeKind is a helper structure to unmarshall the node kind.
type NodeKind struct {
	Kind string `json:"_type"`
}

func unmarshalDeclNode(b []byte) (DeclNode, error) {
	if string(b) == "null" {
		return nil, nil
	}

	var k NodeKind
	if err := json.Unmarshal(b, &k); err != nil {
		return nil, err
	}

	var node DeclNode

	switch nodeKindString[k.Kind] {
	case funcAttr:
		node = &FuncAttr{}
	case funcConfig:
		node = &FuncConfig{}
	case structDecl:
		node = &StructDecl{}
	case enumDecl:
		node = &EnumDecl{}
	default:
		node = &FuncDecl{}
	}

	if err := json.Unmarshal(b, node); err != nil {
		return nil, err
	}
	return node, nil
}

func unmarshalFieldNode(b []byte) (FieldNode, error) {
	if string(b) == "null" {
		return nil, nil
	}

	var k NodeKind
	if err := json.Unmarshal(b, &k); err != nil {
		return nil, err
	}

	var node FieldNode

	switch nodeKindString[k.Kind] {
	case alternative:
		node = &Alternative{}
	default:
		node = &Field{}
	}

	if err := json.Unmarshal(b, node); err != nil {
		return nil, err
	}
	return node, nil
}

func unmarshalTypeNode(b []byte) (TypeNode, error) {
	if string(b) == "null" {
		return nil, nil
	}

	var k NodeKind
	if err := json.Unmarshal(b, &k); err != nil {
		return nil, err
	}

	var node TypeNode

	switch nodeKindString[k.Kind] {
	case dataType:
		node = &DataType{}
	case enumBoolType:
		node = &EnumBoolType{}
	case pointerType:
		node = &PointerType{}
	case sliceType:
		node = &SliceType{}
	case ellipsisType:
		node = &EllipsisType{}
	default:
		node = &BaseType{}
	}

	if err := json.Unmarshal(b, node); err != nil {
		return nil, err
	}
	return node, nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (n *DeclNodes) UnmarshalJSON(b []byte) error {
	*n = DeclNodes{}

	var nodes map[string]json.RawMessage

	if err := json.Unmarshal(b, &nodes); err != nil {
		return err
	}

	for key, raw := range nodes {
		node, err := unmarshalDeclNode(raw)
		if err != nil {
			return err
		}
		(*n)[key] = node
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (l *DeclNodeList) UnmarshalJSON(b []byte) error {
	*l = DeclNodeList{}

	var nodes []json.RawMessage

	if err := json.Unmarshal(b, &nodes); err != nil {
		return err
	}

	for _, raw := range nodes {
		node, err := unmarshalDeclNode(raw)
		if err != nil {
			return err
		}
		*l = append(*l, node)
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (l *FieldList) UnmarshalJSON(b []byte) error {
	*l = FieldList{}

	var nodes []json.RawMessage

	if err := json.Unmarshal(b, &nodes); err != nil {
		return err
	}

	for _, raw := range nodes {
		node, err := unmarshalFieldNode(raw)
		if err != nil {
			return err
		}
		*l = append(*l, node)
	}
	return nil
}

// RawAlternative is a helper structure to unmarshal an alternative and it's type in raw json format.
type RawAlternative struct {
	Name string          `json:"name"`
	Cmd  string          `json:"cmd"`
	Type json.RawMessage `json:"type"`
	List FieldList       `json:"list"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (alt *Alternative) UnmarshalJSON(b []byte) error {
	var ralt RawAlternative
	if err := json.Unmarshal(b, &ralt); err != nil {
		return err
	}
	alt.Name = ralt.Name
	alt.Cmd = ralt.Cmd
	if ralt.Type != nil {
		node, err := unmarshalTypeNode(ralt.Type)
		if err != nil {
			return err
		}
		alt.Type = node
	}
	alt.List = ralt.List
	return nil
}

// RawField is a helper structure to unmarshal a field and it's type in raw json format.
type RawField struct {
	Name string          `json:"name"`
	Cmd  string          `json:"cmd"`
	Type json.RawMessage `json:"type"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (f *Field) UnmarshalJSON(b []byte) error {
	var rf RawField
	if err := json.Unmarshal(b, &rf); err != nil {
		return err
	}
	f.Name = rf.Name
	f.Cmd = rf.Cmd

	if rf.Type != nil {
		node, err := unmarshalTypeNode(rf.Type)
		if err != nil {
			return err
		}
		f.Type = node
	}
	return nil
}

// RawPointerType is a helper structure to unmarshal a pointer.
type RawPointerType struct {
	Node json.RawMessage `json:"node"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *PointerType) UnmarshalJSON(b []byte) error {
	var rt RawPointerType
	if err := json.Unmarshal(b, &rt); err != nil {
		return err
	}
	node, err := unmarshalTypeNode(rt.Node)
	if err != nil {
		return err
	}
	t.Node = node
	return nil
}

// RawSliceType is a helper structure to unmarshal a slice.
type RawSliceType struct {
	AllowNil bool            `json:"allowNil"`
	Variadic bool            `json:"variadic"`
	Node     json.RawMessage `json:"node"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *SliceType) UnmarshalJSON(b []byte) error {
	var rt RawSliceType
	if err := json.Unmarshal(b, &rt); err != nil {
		return err
	}
	t.AllowNil = rt.AllowNil
	t.Variadic = rt.Variadic
	node, err := unmarshalTypeNode(rt.Node)
	if err != nil {
		return err
	}
	t.Node = node
	return nil
}

// RawEllipsisType is a helper structure to unmarshal an ellipsis.
type RawEllipsisType struct {
	Node json.RawMessage `json:"node"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *EllipsisType) UnmarshalJSON(b []byte) error {
	var rt RawEllipsisType
	if err := json.Unmarshal(b, &rt); err != nil {
		return err
	}
	node, err := unmarshalTypeNode(rt.Node)
	if err != nil {
		return err
	}
	t.Node = node
	return nil
}
