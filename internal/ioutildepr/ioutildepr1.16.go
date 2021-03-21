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

// ReadFile reads the named file and returns the contents.
func ReadFile(filename string) ([]byte, error) { return os.ReadFile(filename) }

// WriteFile writes data to a file named by filename.
func WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// ReadAll reads from r until an error or EOF and returns the data it read.
func ReadAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }
