// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package ast

import (
	"sort"
)

// DeclNodes is a map of declaration nodes with key name.
type DeclNodes map[string]DeclNode

// DeclNodeList is a sclice ao all declaration nodes.
type DeclNodeList []DeclNode

// Scope represents an AST scope.
type Scope struct {
	nodes   DeclNodes
	idx     DeclNodeList
	changed bool
}

// NewScope is the Scope constructor.
func NewScope(list DeclNodeList) *Scope {
	s := &Scope{nodes: DeclNodes{}}
	if list != nil {
		for _, decl := range list {
			s.InsertDecl(decl)
		}
	}
	return s
}

// Nodes returns a map of declaration nodes.
func (s *Scope) Nodes() DeclNodes {
	return s.nodes
}

//NodeList return all declaration nodes.
func (s *Scope) NodeList() DeclNodeList {
	// return sorted list
	s.reindex()
	return s.idx
}

// Lookup searches and returns a declaration by name.
func (s *Scope) Lookup(name string) DeclNode {
	return s.nodes[name]
}

// LookupFuncAttr searches and returns a function attribute declaration by name.
func (s *Scope) LookupFuncAttr(name string) *FuncAttr {
	if obj, ok := s.nodes["&"+name]; ok {
		if attr, ok := obj.(*FuncAttr); ok {
			return attr
		}
		panic("invalid function attribute type")
	}
	return nil
}

// LookupFuncConfig searches and returns a function configuration declaration by name.
func (s *Scope) LookupFuncConfig(name string) *FuncConfig {
	if obj, ok := s.nodes["+"+name]; ok {
		if config, ok := obj.(*FuncConfig); ok {
			return config
		}
		panic("invalid function configuration type")
	}
	return nil
}

// InsertDecl inserts a declaration node in scope.
func (s *Scope) InsertDecl(decl DeclNode) bool {
	if _, ok := s.nodes[decl.name()]; ok {
		return false
	}
	s.changed = true
	s.nodes[decl.name()] = decl
	return true
}

func (s *Scope) reindex() {
	if !s.changed {
		return
	}

	s.idx = make([]DeclNode, 0, len(s.nodes))
	for _, decl := range s.nodes {
		s.idx = append(s.idx, decl)
	}
	sort.Slice(s.idx, func(i, j int) bool { return s.idx[i].name() < s.idx[j].name() })

	s.changed = false
}

// LoopFunc iterates through all function declarations.
func (s *Scope) LoopFunc(f func(*FuncDecl)) {
	s.reindex()

	for _, decl := range s.idx {
		if v, ok := decl.(*FuncDecl); ok && !v.Skip {
			f(v)
		}
	}
}

// LoopStruct iterates through all structure declarations.
func (s *Scope) LoopStruct(f func(*StructDecl)) {
	s.reindex()

	for _, decl := range s.idx {
		if v, ok := decl.(*StructDecl); ok {
			f(v)
		}
	}
}

// LoopEnum iterates through all enumeration declarations.
func (s *Scope) LoopEnum(f func(*EnumDecl)) {
	s.reindex()

	for _, decl := range s.idx {
		if v, ok := decl.(*EnumDecl); ok {
			f(v)
		}
	}
}

// DeclNode represents a declaration.
type DeclNode interface {
	declNode()
	name() string
}

// FuncDecl represents a function declaration.
type FuncDecl struct {
	Name  string    `json:"name"`
	Skip  bool      `json:"skip"`
	Attr  string    `json:"attr"`
	Token []string  `json:"token"`
	List  FieldList `json:"list"`
}

// NewFuncDecl is the FuncDecl constructor.
func NewFuncDecl(name, attr string, token []string) *FuncDecl {
	return &FuncDecl{
		Name:  name,
		Attr:  attr,
		Token: token,
		List:  FieldList{},
	}
}

// FuncAttr represents a function attribute declaration.
type FuncAttr struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	Complexity string `json:"complexity"`
	Since      string `json:"since"`
	Group      string `json:"group"`
}

// NewFuncAttr is the FuncAttr constructor.
func NewFuncAttr(name, summary, complexity, since, group string) *FuncAttr {
	return &FuncAttr{
		Name:       name,
		Summary:    summary,
		Complexity: complexity,
		Since:      since,
		Group:      group,
	}
}

// Function configuration constants.
const (
	ConfigType            = "type"
	ConfigTypeSubscribe   = "subscribe"
	ConfigTypeUnsubscribe = "unsubscribe"
	ConfigCallback        = "callback"
	ConfigChannel         = "channel"
)

// FuncConfig represents a function configuration declaration
type FuncConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

// StructDecl represent a structure declaration.
type StructDecl struct {
	Name string    `json:"name"`
	List FieldList `json:"list"`
}

// NewStructDecl is the StructDecl constructor.
func NewStructDecl(name string) *StructDecl {
	return &StructDecl{
		Name: name,
		List: FieldList{},
	}
}

// EnumDecl represents an enumeration declaration.
type EnumDecl struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// NewEnumDecl is the EnumDecl constructor.
func NewEnumDecl(name string, values []string) *EnumDecl {
	return &EnumDecl{
		Name:   name,
		Values: values[:],
	}
}

// DeclNode marker methods
func (t FuncDecl) declNode()   {}
func (t FuncAttr) declNode()   {}
func (t FuncConfig) declNode() {}
func (t StructDecl) declNode() {}
func (t EnumDecl) declNode()   {}

// DeclNode name methods
func (t *FuncDecl) name() string   { return t.Name }
func (t *FuncAttr) name() string   { return "&" + t.Name }
func (t *FuncConfig) name() string { return "+" + t.Name }
func (t *StructDecl) name() string { return t.Name }
func (t *EnumDecl) name() string   { return t.Name }

// check if type implements DeclNode interface
var _ DeclNode = (*FuncDecl)(nil)
var _ DeclNode = (*FuncAttr)(nil)
var _ DeclNode = (*FuncConfig)(nil)
var _ DeclNode = (*StructDecl)(nil)
var _ DeclNode = (*EnumDecl)(nil)

// FieldNode represents a field definition.
type FieldNode interface {
	fieldNode()
	NodeName() string
	NodeCmd() string
	NodeType() TypeNode
	walkNode(level int, fct WalkNodeFct)
}

// WalkNodeFct is the callback function definition for node visitor WalkNode.
type WalkNodeFct func(level int, node FieldNode)

// FieldList is a list of field definitions.
type FieldList []FieldNode

// WalkNode iterates through all nodes of the list.
func (list FieldList) WalkNode(fct WalkNodeFct) {
	list.walkNode(0, fct)
}

func (list FieldList) walkNode(level int, fct WalkNodeFct) {
	for _, node := range list {
		node.walkNode(level, fct)
	}
}

// Alternative represents a set of alternative fields.
type Alternative struct {
	Name string    `json:"name"`
	Cmd  string    `json:"cmd"`
	Type TypeNode  `json:"type"`
	List FieldList `json:"list"`
}

// NodeName implements FieldNode interface.
func (alt *Alternative) NodeName() string {
	return alt.Name
}

// NodeCmd implements FieldNode interface.
func (alt *Alternative) NodeCmd() string {
	return alt.Cmd
}

// NodeType implements FieldNode interface.
func (alt *Alternative) NodeType() TypeNode {
	return alt.Type
}

func (alt *Alternative) walkNode(level int, fct WalkNodeFct) {
	fct(level, alt)
	alt.List.walkNode(level+1, fct)
}

// Field represents an AST field.
type Field struct {
	Name string   `json:"name"`
	Cmd  string   `json:"cmd"`
	Type TypeNode `json:"type"`
}

// NodeName implements FieldNode interface.
func (f *Field) NodeName() string {
	return f.Name
}

// NodeCmd implements FieldNode interface.
func (f *Field) NodeCmd() string {
	return f.Cmd
}

// NodeType implements FieldNode interface.
func (f *Field) NodeType() TypeNode {
	return f.Type
}

func (f *Field) walkNode(level int, fct WalkNodeFct) {
	fct(level, f)
}

// FieldNode marker methods
func (alt Alternative) fieldNode() {}
func (f Field) fieldNode()         {}

// check if type implements FieldNode interface
var _ FieldNode = (*Alternative)(nil)
var _ FieldNode = (*Field)(nil)

// TypeNode represents a node with type information.
type TypeNode interface {
	typeNode()
	String() string
}

var (
	// KeyType represents a key.
	KeyType TypeNode = &BaseType{Name: "interface{}"}
	// InterfaceType represents a generic type.
	InterfaceType TypeNode = &BaseType{Name: "interface{}"}
	// StringType represents a string.
	StringType TypeNode = &BaseType{Name: "string"}
	// IntegerType represents an integer.
	IntegerType TypeNode = &BaseType{Name: "int64"}
	// FloatType represents a float.
	FloatType TypeNode = &BaseType{Name: "float64"}
	// BoolType represents a boolean.
	BoolType TypeNode = &BaseType{Name: "bool"}
	// TimeType represents a time.
	TimeType TypeNode = &BaseType{Name: "time.Time"}

	// StringPointerType represents a pointer to a string
	StringPointerType TypeNode = &PointerType{Node: StringType}
	// IntegerPointerType represents a pointer to an integer.
	IntegerPointerType TypeNode = &PointerType{Node: IntegerType}

	// StringSliceType represents a slice of string.
	StringSliceType TypeNode = &SliceType{Node: StringType}
)

// BaseType represents a base type (like string, int64, ...).
type BaseType struct {
	Name string `json:"name"`
}

// DataType represents a structure.
type DataType struct {
	Name string `json:"name"`
}

// EnumBoolType represents an one or two element enumeration.
type EnumBoolType struct {
	Values []string `json:"values"`
}

// BaseType returns the base type of an EnumBoolType.
func (t *EnumBoolType) BaseType() TypeNode {
	return BoolType
}

// PointerType represents an optinal attribute.
type PointerType struct {
	Node TypeNode `json:"node"`
}

// SliceType represents a multi-value attribute.
type SliceType struct {
	AllowNil bool     `json:"allowNil"`
	Variadic bool     `json:"variadic"`
	Node     TypeNode `json:"node"`
}

// EllipsisType represents a variadic attribute.
type EllipsisType struct {
	Node TypeNode `json:"node"`
}

// TypeNode marker methods
func (t BaseType) typeNode()     {}
func (t DataType) typeNode()     {}
func (t EnumBoolType) typeNode() {}
func (t PointerType) typeNode()  {}
func (t SliceType) typeNode()    {}
func (t EllipsisType) typeNode() {}

// TypeNode Stringer interface
func (t *BaseType) String() string     { return t.Name }
func (t *DataType) String() string     { return t.Name }
func (t *EnumBoolType) String() string { return BoolType.String() }
func (t *PointerType) String() string  { return "*" + t.Node.String() }
func (t *SliceType) String() string    { return "[]" + t.Node.String() }
func (t *EllipsisType) String() string { return "..." + t.Node.String() }

// check if type implements TypeNode interface
var _ TypeNode = (*BaseType)(nil)
var _ TypeNode = (*DataType)(nil)
var _ TypeNode = (*EnumBoolType)(nil)
var _ TypeNode = (*PointerType)(nil)
var _ TypeNode = (*SliceType)(nil)
var _ TypeNode = (*EllipsisType)(nil)
