// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"math/big"
	"strconv"
)

// RedisKind represents the kind of type that a RedisValue represents.
type RedisKind byte

// RedisKind constants.
const (
	RkInvalid RedisKind = iota
	RkError
	RkPush
	RkNull
	RkString
	RkVerbatimString
	RkNumber
	RkDouble
	RkBigNumber
	RkBoolean
	RkSlice
	RkMap
	RkSet
)

// RedisValue represents a Redis command reply value.
type RedisValue interface {
	Kind() RedisKind
	Attr() *Map
	Converter
}

type baseRedisType interface {
	_interface() interface{}
}

var (
	_Null           = _null{}
	_VerbatimString = _verbatimString("")
	_Slice          = _slice{}
	_Map            = _map{}
	_Set            = _set{}
)

var _ baseRedisType = (*_null)(nil)
var _ baseRedisType = (*_string)(nil)
var _ baseRedisType = (*_number)(nil)
var _ baseRedisType = (*_bignumber)(nil)
var _ baseRedisType = (*_double)(nil)
var _ baseRedisType = (*_boolean)(nil)

var _ RedisValue = (*_null)(nil)
var _ RedisValue = (*_string)(nil)
var _ RedisValue = (*_number)(nil)
var _ RedisValue = (*_bignumber)(nil)
var _ RedisValue = (*_double)(nil)
var _ RedisValue = (*_boolean)(nil)

var _ RedisValue = (*_verbatimString)(nil)
var _ RedisValue = (*_slice)(nil)
var _ RedisValue = (*_map)(nil)
var _ RedisValue = (*_set)(nil)

var _ RedisValue = (*attrRedisValue)(nil)

type attrRedisValue struct {
	RedisValue
	attr _map
}

func (v attrRedisValue) Attr() *Map { return (*Map)(&v.attr) }

// Caution: _null need to implement all Slice, Map, Set conversion functions (support method chaining)
type _null struct{}

func (n _null) _interface() interface{} { return nil }
func (n _null) Kind() RedisKind         { return RkNull }

func (n _null) ToSlice() (Slice, error)                             { return _Slice.ToSlice() }
func (n _null) ToInt64Slice() ([]int64, error)                      { return _Slice.ToInt64Slice() }
func (n _null) ToIntfSlice() ([]interface{}, error)                 { return _Slice.ToIntfSlice() }
func (n _null) ToIntfSlice2() ([][]interface{}, error)              { return _Slice.ToIntfSlice2() }
func (n _null) ToIntfSlice3() ([][][]interface{}, error)            { return _Slice.ToIntfSlice3() }
func (n _null) ToStringMapSlice() ([]map[string]interface{}, error) { return _Slice.ToStringMapSlice() }
func (n _null) ToStringSlice() ([]string, error)                    { return _Slice.ToStringSlice() }
func (n _null) ToTree() ([]interface{}, error)                      { return _Slice.ToTree() }
func (n _null) ToXrange() ([]XItem, error)                          { return _Slice.ToXrange() }
func (n _null) ToMap() (Map, error)                                 { return _Map.ToMap() }
func (n _null) ToStringInt64Map() (map[string]int64, error)         { return _Map.ToStringInt64Map() }
func (n _null) ToStringMap() (map[string]interface{}, error)        { return _Map.ToStringMap() }
func (n _null) ToStringValueMap() (map[string]RedisValue, error)    { return _Map.ToStringValueMap() }
func (n _null) ToStringStringMap() (map[string]string, error)       { return _Map.ToStringStringMap() }
func (n _null) ToXread() (map[string][]XItem, error)                { return _Map.ToXread() }
func (n _null) ToSet() (Set, error)                                 { return _Set.ToSet() }
func (n _null) ToStringSet() (map[string]bool, error)               { return _Set.ToStringSet() }

type _string string

func (s _string) _interface() interface{}     { return string(s) }
func (s _string) Kind() RedisKind             { return RkString }
func (s _string) ToString() (string, error)   { return string(s), nil }
func (s _string) ToInt64() (int64, error)     { return strconv.ParseInt(string(s), 10, 64) }
func (s _string) ToFloat64() (float64, error) { return strconv.ParseFloat(string(s), 64) }
func (s _string) ToBool() (bool, error)       { return s == ReplyOK, nil }

type _number int64

func (n _number) _interface() interface{}     { return int64(n) }
func (n _number) Kind() RedisKind             { return RkNumber }
func (n _number) ToString() (string, error)   { return strconv.FormatInt(int64(n), 10), nil }
func (n _number) ToInt64() (int64, error)     { return int64(n), nil }
func (n _number) ToFloat64() (float64, error) { return float64(n), nil }
func (n _number) ToBool() (bool, error)       { return n != 0, nil }

type _double float64

func (d _double) _interface() interface{}     { return float64(d) }
func (d _double) Kind() RedisKind             { return RkDouble }
func (d _double) ToString() (string, error)   { return strconv.FormatFloat(float64(d), 'g', -1, 64), nil }
func (d _double) ToFloat64() (float64, error) { return float64(d), nil }
func (d _double) ToBool() (bool, error)       { return d != 0, nil }

type _bignumber big.Int

func (n *_bignumber) _interface() interface{}   { return (*big.Int)(n) }
func (n *_bignumber) Kind() RedisKind           { return RkBigNumber }
func (n *_bignumber) ToString() (string, error) { return (*big.Int)(n).String(), nil }
func (n *_bignumber) ToInt64() (int64, error) {
	if (*big.Int)(n).IsInt64() {
		return (*big.Int)(n).Int64(), nil
	}
	return 0, newConversionError("ToInt64", n)
}
func (n *_bignumber) ToFloat64() (float64, error) {
	if (*big.Int)(n).IsInt64() {
		return float64((*big.Int)(n).Int64()), nil
	}
	return 0, newConversionError("ToFloat64", n)
}
func (n *_bignumber) ToBool() (bool, error) { return (*big.Int)(n).Sign() != 0, nil }

type _boolean bool

func (b _boolean) _interface() interface{}   { return bool(b) }
func (b _boolean) Kind() RedisKind           { return RkBoolean }
func (b _boolean) ToString() (string, error) { return strconv.FormatBool(bool(b)), nil }
func (b _boolean) ToInt64() (int64, error) {
	if b {
		return 1, nil
	}
	return 0, nil
}
func (b _boolean) ToFloat64() (float64, error) {
	if b {
		return 1, nil
	}
	return 0, nil
}
func (b _boolean) ToBool() (bool, error) { return bool(b), nil }

type _verbatimString string

func (s _verbatimString) Kind() RedisKind                           { return RkVerbatimString }
func (s _verbatimString) ToVerbatimString() (VerbatimString, error) { return VerbatimString(s), nil }
func (s _verbatimString) FileFormat() string                        { return string(s[:3]) }
func (s _verbatimString) String() string                            { return string(s[4:]) }
func (s _verbatimString) ToString() (string, error)                 { return string(s[4:]), nil }
func (s _verbatimString) ToInt64() (int64, error)                   { return strconv.ParseInt(string(s[4:]), 10, 64) }
func (s _verbatimString) ToFloat64() (float64, error)               { return strconv.ParseFloat(string(s[4:]), 64) }
func (s _verbatimString) ToBool() (bool, error)                     { return string(s[4:]) == ReplyOK, nil }

type _slice []RedisValue

func (s _slice) Kind() RedisKind         { return RkSlice }
func (s _slice) ToSlice() (Slice, error) { return Slice(s), nil }
func (s _slice) ToIntfSlice() ([]interface{}, error) {
	r := make([]interface{}, len(s))
	for i, item := range s {
		switch value := item.(type) {
		case baseRedisType:
			r[i] = value._interface()
		default:
			r[i] = value
		}
	}
	return r, nil
}
func (s _slice) ToIntfSlice2() ([][]interface{}, error) {
	r := make([][]interface{}, len(s))
	for i, item := range s {
		l, err := item.ToIntfSlice()
		if err != nil {
			return nil, newConversionError("toSlice", item)
		}
		r[i] = l
	}
	return r, nil
}
func (s _slice) ToIntfSlice3() ([][][]interface{}, error) {
	r := make([][][]interface{}, len(s))
	for i, item := range s {
		l, err := item.ToIntfSlice2()
		if err != nil {
			return nil, newConversionError("toSlice2", item)
		}
		r[i] = l
	}
	return r, nil
}
func (s _slice) ToTree() ([]interface{}, error) {
	r := make([]interface{}, len(s))
	for i, item := range s {
		if item.Kind() == RkSlice {
			var err error
			r[i], err = item.ToTree()
			if err != nil {
				return nil, err
			}
		} else {
			switch value := item.(type) {
			case baseRedisType:
				r[i] = value._interface()
			default:
				r[i] = value
			}
		}
	}
	return r, nil
}
func (s _slice) ToStringSlice() ([]string, error) {
	r := make([]string, len(s))
	for i, item := range s {
		val, err := item.ToString()
		if err != nil {
			return nil, err
		}
		r[i] = val
	}
	return r, nil
}
func (s _slice) ToStringMapSlice() ([]map[string]interface{}, error) {
	r := make([]map[string]interface{}, len(s))
	for i, item := range s {
		val, err := item.ToStringMap()
		if err != nil {
			return nil, err
		}
		r[i] = val
	}
	return r, nil
}
func (s _slice) ToInt64Slice() ([]int64, error) {
	r := make([]int64, len(s))
	for i, item := range s {
		val, err := item.ToInt64()
		if err != nil {
			return nil, err
		}
		r[i] = val
	}
	return r, nil
}
func (s _slice) ToXrange() ([]XItem, error) {
	r := make([]XItem, len(s))
	for i, item := range s {
		if item.Kind() != RkSlice {
			return nil, newConversionError("toXrange", item)
		}
		slice := item.(_slice)
		if len(slice) != 2 {
			return nil, newConversionError("toXrange", item)
		}
		s, err := slice[0].ToString()
		if err != nil {
			return nil, err
		}
		l, err := slice[1].ToStringSlice()
		if err != nil {
			return nil, err
		}
		r[i] = XItem{s, l}
	}
	return r, nil
}

type _map []MapItem

func (m _map) Kind() RedisKind     { return RkMap }
func (m _map) ToMap() (Map, error) { return Map(m), nil }
func (m _map) ToStringMap() (map[string]interface{}, error) {
	r := make(map[string]interface{}, len(m))
	for _, item := range m {
		key, err := item.Key.ToString()
		if err != nil {
			return nil, err
		}
		switch value := item.Value.(type) {
		case baseRedisType:
			r[key] = value._interface()
		default:
			r[key] = value
		}
	}
	return r, nil
}
func (m _map) ToStringValueMap() (map[string]RedisValue, error) {
	r := make(map[string]RedisValue, len(m))
	for _, item := range m {
		key, err := item.Key.ToString()
		if err != nil {
			return nil, err
		}
		r[key] = item.Value
	}
	return r, nil
}
func (m _map) ToStringStringMap() (map[string]string, error) {
	r := make(map[string]string, len(m))
	for _, item := range m {
		key, err := item.Key.ToString()
		if err != nil {
			return nil, err
		}
		value, err := item.Value.ToString()
		if err != nil {
			return nil, err
		}
		r[key] = value
	}
	return r, nil
}
func (m _map) ToStringInt64Map() (map[string]int64, error) {
	r := make(map[string]int64, len(m))
	for _, item := range m {
		key, err := item.Key.ToString()
		if err != nil {
			return nil, err
		}
		value, err := item.Value.ToInt64()
		if err != nil {
			return nil, err
		}
		r[key] = value
	}
	return r, nil
}
func (m _map) ToXread() (map[string][]XItem, error) {
	r := make(map[string][]XItem, len(m))
	for _, item := range m {
		s, err := item.Key.ToString()
		if err != nil {
			return nil, err
		}
		xrange, err := item.Value.ToXrange()
		if err != nil {
			return nil, err
		}
		r[s] = xrange
	}
	return r, nil
}

type _set []RedisValue

func (s _set) Kind() RedisKind     { return RkSet }
func (s _set) ToSet() (Set, error) { return Set(s), nil }
func (s _set) ToStringSet() (map[string]bool, error) {
	r := make(map[string]bool, len(s))
	for _, item := range s {
		key, err := item.ToString()
		if err != nil {
			return nil, err
		}
		r[key] = true
	}
	return r, nil
}
