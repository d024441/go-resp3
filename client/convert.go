// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
)

// A ConversionError is raised for an unsuccessful type conversion of a redis value.
// - To:    Name of the conversion function.
// - Value: Value for which the conversion was not successful.
type ConversionError struct {
	To    string
	Value interface{}
}

func newConversionError(to string, value interface{}) *ConversionError {
	return &ConversionError{to, value}
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("unsupported %s conversion type %T", e.To, e.Value)
}

// Converter is the interface that groups Redis value conversion interfaces.
type Converter interface {
	Stringer
	VerbatimStringer
	Int64er
	Float64er
	Booler
	Slicer
	Mapper
	Setter

	Int64Slicer
	IntfSlicer
	IntfSlice2er
	IntfSlice3er
	StringMapSlicer
	StringSlicer
	Treer
	Xranger

	StringMapper
	StringValueMapper
	StringInt64Mapper
	StringStringMapper
	Xreader

	StringSetter
}

// Stringer is implemented by any redis value that has a ToString method.
type Stringer interface {
	// ToString converts a redis value to a string.
	// In case the conversion is not supported a ConversionError is returned.
	ToString() (string, error)
}

// VerbatimStringer is implemented by any redis value that has a ToVerbatimString method.
type VerbatimStringer interface {
	// ToVerbatimString converts a redis value to a VerbatimString.
	// In case the conversion is not supported a ConversionError is returned.
	ToVerbatimString() (VerbatimString, error)
}

// Int64er is implemented by any redis value that has a ToInt64 method.
type Int64er interface {
	// ToInt64 converts a redis value to an int64.
	// In case the conversion is not supported a ConversionError is returned.
	ToInt64() (int64, error)
}

// Float64er is implemented by any redis value that has a ToFloat64 method.
type Float64er interface {
	// ToFloat64 converts a redis value to a float64.
	// In case the conversion is not supported a ConversionError is returned.
	ToFloat64() (float64, error)
}

// Booler is implemented by any redis value that has a ToBool method.
type Booler interface {
	// ToBooler converts a redis value to a bool.
	// In case the conversion is not supported a ConversionError is returned.
	ToBool() (bool, error)
}

// Slicer is implemented by any redis value that has a ToSlice method.
type Slicer interface {
	// ToSlice converts a redis value to a Slice.
	// In case value conversion is not possible a ConversitionError is returned.
	ToSlice() (Slice, error)
}

// Mapper is implemented by any redis value that has a ToMap method.
type Mapper interface {
	// ToMap converts a redis value to a Map.
	// In case value conversion is not possible a ConversitionError is returned.
	ToMap() (Map, error)
}

// Setter is implemented by any redis value that has a ToSet method.
type Setter interface {
	// ToSet converts a redis value to a Set.
	// In case value conversion is not possible a ConversitionError is returned.
	ToSet() (Set, error)
}

// Int64Slicer is implemented by any redis value that has a ToInt64Slice method.
type Int64Slicer interface {
	// ToInt64Slice returns a slice with values of type int64. In case value conversion to []int64 is not possible
	// a ConversitionError is returned.
	ToInt64Slice() ([]int64, error)
}

// IntfSlicer is implemented by any redis value that has a ToIntfSlice method.
type IntfSlicer interface {
	// ToIntfSlice returns a slice with values of type interface{}.
	ToIntfSlice() ([]interface{}, error)
}

// IntfSlice2er is implemented by any redis value that has a ToIntfSlice2 method.
type IntfSlice2er interface {
	// ToInttSlice2 returns a slice with values of type []interface{}. In case value conversion to []interface{} is not possible
	// a ConversitionError is returned.
	ToIntfSlice2() ([][]interface{}, error)
}

// IntfSlice3er is implemented by any redis value that has a ToIntfSlice3 method.
type IntfSlice3er interface {
	// ToIntSlice3 returns a slice with values of type [][]interface{}. In case value conversion to [][]interface{} is not possible
	// a ConversitionError is returned.
	ToIntfSlice3() ([][][]interface{}, error)
}

// StringMapSlicer is implemented by any redis value that has a ToStringMapSlice method.
type StringMapSlicer interface {
	// ToStringMapSlice returns a slice with values of type map[string]interface{}. In case value conversion to map[string]interface{} is not possible
	// a ConversitionError is returned.
	ToStringMapSlice() ([]map[string]interface{}, error)
}

// StringSlicer is implemented by any redis value that has a ToStringSlice method.
type StringSlicer interface {
	// ToStringSlice returns a slice with values of type string. In case value conversion to []string is not possible
	// a ConversitionError is returned.
	ToStringSlice() ([]string, error)
}

// Treer is implemented by any redis value that has a ToTree method.
type Treer interface {
	// ToTree returns a tree with nodes of type []interface{} and leaves of type interface{}. In case value conversion to []interface{} is not possible
	// a ConversitionError is returned.
	ToTree() ([]interface{}, error)
}

// Xranger is implemented by any redis value that has a ToXrange method.
type Xranger interface {
	// ToXrange returns a slice with values of type XItem. In case the conversion is not possible
	// a ConversitionError is returned.
	ToXrange() ([]XItem, error)
}

// StringInt64Mapper is implemented by any redis value that has a ToStringInt64Map method.
type StringInt64Mapper interface {
	// ToStringInt64Map returns a map with keys of type string and values of type int64. In case key or value conversion is not possible
	// a ConvertionError is returned.
	ToStringInt64Map() (map[string]int64, error)
}

// StringMapper is implemented by any redis value that has a ToStringMap method.
type StringMapper interface {
	// ToStringMap returns a map with keys of type string. In case key conversion to string is not possible
	// a ConvertionError is returned.
	ToStringMap() (map[string]interface{}, error)
}

// StringValueMapper is implemented by any redis value that has a ToStringValueMap method.
type StringValueMapper interface {
	// ToStringValueMap returns a map with keys of type string and values of type RedisValue.
	// In case key conversion to string is not possible a ConvertionError is returned.
	ToStringValueMap() (map[string]RedisValue, error)
}

// StringStringMapper is implemented by any redis value that has a ToStringStringMap method.
type StringStringMapper interface {
	// ToStringStringMap returns a map with keys and values of type string. In case key or value conversion to string is not possible
	// a ConvertionError is returned.
	ToStringStringMap() (map[string]string, error)
}

// Xreader is implemented by any redis value that has a ToXread method.
type Xreader interface {
	// ToXread returns a map[string] with values of type XItem. In case the conversion is not possible
	// a ConversitionError is returned.
	ToXread() (map[string][]XItem, error)
}

// StringSetter is implemented by any redis value that has a ToStringSet method.
type StringSetter interface {
	// ToStringSet returns a map with keys of type string and boolean true values. In case key conversion to string is not possible
	// a ConvertionError is returned.
	ToStringSet() (map[string]bool, error)
}
