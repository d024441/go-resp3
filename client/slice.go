// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

// A Slice represents the redis slice type.
type Slice []RedisValue

// ToIntfSlice returns a slice with values of type interface{}.
func (s Slice) ToIntfSlice() ([]interface{}, error) { return _slice(s).ToIntfSlice() }

// ToIntfSlice2 returns a slice with values of type []interface{}. In case value conversion to []interface{} is not possible
// a ConversitionError is returned.
func (s Slice) ToIntfSlice2() ([][]interface{}, error) { return _slice(s).ToIntfSlice2() }

// ToIntfSlice3 returns a slice with values of type [][]interface{}. In case value conversion to [][]interface{} is not possible
// a ConversitionError is returned.
func (s Slice) ToIntfSlice3() ([][][]interface{}, error) { return _slice(s).ToIntfSlice3() }

// ToTree returns a tree with nodes of type []interface{} and leaves of type interface{}. In case value conversion to []interface{} is not possible
// a ConversitionError is returned.
func (s Slice) ToTree() ([]interface{}, error) { return _slice(s).ToTree() }

// ToStringSlice returns a slice with values of type string. In case value conversion to string is not possible
// a ConversitionError is returned.
func (s Slice) ToStringSlice() ([]string, error) { return _slice(s).ToStringSlice() }

// ToStringMapSlice returns a slice with values of type map[string]interface{}. In case value conversion to map[string]interface{} is not possible
// a ConversitionError is returned.
func (s Slice) ToStringMapSlice() ([]map[string]interface{}, error) {
	return _slice(s).ToStringMapSlice()
}

// ToInt64Slice returns a slice with values of type int64. In case value conversion to string is not possible
// a ConversitionError is returned.
func (s Slice) ToInt64Slice() ([]int64, error) { return _slice(s).ToInt64Slice() }

// ToXrange returns a slice with values of type XItem. In case the conversion is not possible
// a ConversitionError is returned.
func (s Slice) ToXrange() ([]XItem, error) { return _slice(s).ToXrange() }
