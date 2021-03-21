// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package conv

import (
	. "testing"
)

func TestParseInt(t *T) {
	var parseTest = []struct {
		s string
		n int64
	}{
		{"", 0}, {"1", 1}, {"+1", 1}, {"-1", -1}, {"99999", 99999},
	}

	for i, test := range parseTest {
		n, err := ParseInt([]byte(test.s))
		if err != nil {
			t.Fatal(err)
		}
		if n != test.n {
			t.Fatalf("%d: %d - %d", i, n, test.n)
		}
	}
}
