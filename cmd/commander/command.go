// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
)

type commands map[string]*command

type command struct {
	Summary    string      `json:"summary"`
	Complexity string      `json:"complexity"`
	Arguments  []*argument `json:"arguments"`
	Since      string      `json:"since"`
	Group      string      `json:"group"`
}

type vName struct {
	id  string
	ids []string
}

func (v *vName) UnmarshalJSON(b []byte) error {
	if b[0] == '[' { // array
		err := json.Unmarshal(b, &v.ids)
		if len(v.ids) == 1 { // correct def if only one element in array
			v.id = v.ids[0]
			v.ids = nil
		}
		return err
	}
	return json.Unmarshal(b, &v.id)
}

type vType struct {
	id  string
	ids []string
}

func (v *vType) UnmarshalJSON(b []byte) error {
	if b[0] == '[' { // array
		err := json.Unmarshal(b, &v.ids)
		if len(v.ids) == 1 { // correct def if only one element in array
			v.id = v.ids[0]
			v.ids = nil
		}
		return err
	}
	return json.Unmarshal(b, &v.id)
}

type argument struct {
	Command  string   `json:"command"`
	Name     vName    `json:"name"`
	Type     vType    `json:"type"`
	Enum     []string `json:"enum"`
	Optional bool     `json:"optional"`
	Multiple bool     `json:"multiple"` // multiple args (including multiple command)
	Variadic bool     `json:"variadic"` // variadic args (command only once)
}

type argKind int

const (
	akUnknown   = iota
	akEnumConst // enum used as text constant
	akEnumBool  // enum can be represented as bool
	akEnum      // normal enum type
	akStruct    // struct type
	akBasic     // basic type
)

func (a *argument) Kind() argKind {
	if a.Enum != nil {
		switch len(a.Enum) {
		case 1:
			if a.Optional {
				return akEnumBool
			}
			return akEnumConst
		case 2:
			return akEnumBool
		default:
			return akEnum
		}
	}
	if a.Name.ids != nil {
		return akStruct
	}
	return akBasic
}
