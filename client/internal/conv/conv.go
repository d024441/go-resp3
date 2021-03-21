// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package conv

import (
	"errors"
)

// ErrInvalidCharacter is returned if parse functions detect an invalid character.
var ErrInvalidCharacter = errors.New("invalid character")

// ParseInt parses an integer.
// - Like strings ParseInt but avoiding string allocation.
func ParseInt(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, nil
	}
	switch b[0] {
	case '+':
		n, err := ParseUint(b[1:])
		return int64(n), err
	case '-':
		n, err := ParseUint(b[1:])
		return int64(n) * -1, err
	default:
		n, err := ParseUint(b)
		return int64(n), err
	}
}

// ParseUint parses an unsigned integer.
// - Like strings ParseUint but avoiding string allocation.
func ParseUint(b []byte) (uint64, error) {
	var n uint64
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, ErrInvalidCharacter
		}
		d := c - '0'
		n = n*10 + uint64(d)
	}
	return n, nil
}
