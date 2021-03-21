// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package stringutils

import (
	"strings"
	"unicode"
)

// PascalCase returns the pascal case version of a string.
func PascalCase(s string) string {
	return multiCase(true, s)
}

// CamelCase returns the camel case version of a string.
func CamelCase(s string) string {
	return multiCase(false, s)
}

func multiCase(capital bool, s string) string {
	j := 0
	result := make([]rune, len(s))
	for _, r := range s {
		switch r {
		case ' ', '-', '_', ':', '/':
			capital = true
		default:
			if capital {
				capital = false
				result[j] = unicode.ToUpper(r)

			} else {
				result[j] = unicode.ToLower(r)
			}
			j++
		}
	}
	return string(result[:j])
}

// Split splits a string in lines, whare each line does not have more than max characters.
func Split(s string, max int) []string {
	r := make([]string, 0)
	parts := strings.SplitAfter(s, " ")

	l, i := 0, 0
	for j, part := range parts {
		l += len(part)
		if l > max {
			r = append(r, strings.Join(parts[i:j], ""))
			i, l = j, 0
		}
	}
	if i < len(parts) {
		r = append(r, strings.Join(parts[i:], ""))
	}
	return r
}
