// +build debug

// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"reflect"
)

func reflectStruct(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("Field %s size %d align %d offset %d\n", f.Name, f.Type.Size(), f.Type.Align(), f.Offset)
	}
}

func reflectType(t reflect.Type) {
	switch t.Kind() {

	case reflect.Struct:
		fmt.Printf("Struct %s size %d align %d\n", t.Name(), t.Size(), t.Align())
		reflectStruct(t)
	}
}

func init() {
	reflectType(reflect.TypeOf(result{}))
	reflectType(reflect.TypeOf(request{}))
}
