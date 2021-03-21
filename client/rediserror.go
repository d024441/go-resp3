// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"strings"
)

// A RedisError represents the redis error message if a redis command was executed unsuccessfully.
type RedisError struct {
	Code string
	Msg  string
}

func newRedisError(s string) *RedisError {
	p := strings.Split(s, " ")
	return &RedisError{
		Code: p[0],
		Msg:  strings.Join(p[1:], " "),
	}
}

func (e *RedisError) Error() string { return e.Code + " " + e.Msg }
