// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// RandomKey returns key + base64 encoded random bytes. Used for tests to avoid key overwrites on Redis server.
func RandomKey(key string) string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return key + base64.URLEncoding.EncodeToString(b)
}
