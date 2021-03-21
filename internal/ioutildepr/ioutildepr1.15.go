// +build !go1.16

// SPDX-FileCopyrightText: 2014-2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Delete after go1.15 is out of maintenance.

package ioutildepr

import (
	"io"
	"io/ioutil"
	"os"
)

func ReadFile(filename string) ([]byte, error) { return ioutil.ReadFile(filename) }
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
func ReadAll(r io.Reader) ([]byte, error) { return ioutil.ReadAll(r) }
