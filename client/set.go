// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

// A Set represents the redis set type.
type Set []RedisValue

// ToStringSet returns a map with keys of type string and boolean true values. In case key conversion to string is not possible
// a ConvertionError is returned.
func (s Set) ToStringSet() (map[string]bool, error) { return _set(s).ToStringSet() }
