// +build go1.16

// SPDX-FileCopyrightText: 2014-2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Delete and re-itegrate after go1.15 is out of maintenance.

package ioutildepr

import (
	"io"
	"io/fs"
	"os"
)

func ReadFile(filename string) ([]byte, error) { return os.ReadFile(filename) }
func WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
func ReadAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }
