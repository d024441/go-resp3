// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

// A VerbatimString represents the redis verbatim string type.
type VerbatimString string

// FileFormat is returning the file format of a verbatim string.
func (s VerbatimString) FileFormat() string { return _verbatimString(s).FileFormat() }

// String implements the Stringer interface.
func (s VerbatimString) String() string { return _verbatimString(s).String() }

// ToString implements the Converter Stringer interface.
func (s VerbatimString) ToString() (string, error) { return _verbatimString(s).ToString() }

// ToInt64 implements the Converter ToInt64er interface.
func (s VerbatimString) ToInt64() (int64, error) { return _verbatimString(s).ToInt64() }

// ToFloat64 implements the Converter ToFloat64er interface.
func (s VerbatimString) ToFloat64() (float64, error) { return _verbatimString(s).ToFloat64() }

// ToBool implements the Converter ToBooler interface.
func (s VerbatimString) ToBool() (bool, error) { return _verbatimString(s).ToBool() }
