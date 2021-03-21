// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

// A MapItem represents the redis map type key and value.
type MapItem struct {
	Key   RedisValue
	Value RedisValue
}

// A Map represents the redis map type.
type Map []MapItem

// ToStringMap returns a map with keys of type string. In case key conversion to string is not possible
// a ConvertionError is returned.
func (m Map) ToStringMap() (map[string]interface{}, error) { return _map(m).ToStringMap() }

// ToStringValueMap returns a map with keys of type stringand values of type RedisValue.
// In case key conversion to string is not possible a ConvertionError is returned.
func (m Map) ToStringValueMap() (map[string]RedisValue, error) { return _map(m).ToStringValueMap() }

// ToStringStringMap returns a map with keys and values of type string. In case key or value conversion to string is not possible
// a ConvertionError is returned.
func (m Map) ToStringStringMap() (map[string]string, error) { return _map(m).ToStringStringMap() }

// ToStringInt64Map returns a map with keys of type string and values of type int64. In case key or value conversion is not possible
// a ConvertionError is returned.
func (m Map) ToStringInt64Map() (map[string]int64, error) { return _map(m).ToStringInt64Map() }

// XItem represents the id and the entry list of streams.
type XItem struct {
	ID    string
	Items []string
}

// ToXread returns a map[string] with values of type XItem. In case the conversion is not possible
// a ConversitionError is returned.
func (m Map) ToXread() (map[string][]XItem, error) { return _map(m).ToXread() }
