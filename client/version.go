// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Version holds the information of a semantic version (https://semver.org/).
type Version struct {
	Major, Minor, Patch       uint64
	Prerelease, Buildmetadata string
}

// ParseVersion parses a semantic version string field.
func ParseVersion(s string) Version {

	v := Version{}
	parts := strings.SplitN(s, ".", 3)
	l := len(parts)
	if l < 1 {
		return v
	}
	v.Major, _ = strconv.ParseUint(parts[0], 10, 64)
	if l < 2 {
		return v
	}
	v.Minor, _ = strconv.ParseUint(parts[1], 10, 64)
	if l < 3 {
		return v
	}
	parts = strings.SplitN(parts[2], "+", 2)
	if len(parts) > 1 {
		v.Buildmetadata = parts[1]
	}
	parts = strings.SplitN(parts[0], "-", 2)
	if len(parts) > 1 {
		v.Prerelease = parts[1]
	}
	v.Patch, _ = strconv.ParseUint(parts[0], 10, 64)
	return v
}

func (v Version) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		fmt.Fprintf(&b, "-%s", v.Prerelease)
	}
	if v.Buildmetadata != "" {
		fmt.Fprintf(&b, "+%s", v.Buildmetadata)
	}
	return b.String()
}

func compareUint64(u1, u2 uint64) int {
	switch {
	case u1 == u2:
		return 0
	case u1 > u2:
		return 1
	default:
		return -1
	}
}

func minInt(l1, l2 int) int {
	if l1 < l2 {
		return l1
	}
	return l2
}

// Compare compares the version with a second version v2. The result will be
//  0 in case the two versions are equal,
// -1 in case version v has lower precedence than c2,
//  1 in case version v has higher precedence than c2.
func (v Version) Compare(v2 Version) int {
	if r := compareUint64(v.Major, v2.Major); r != 0 {
		return r
	}
	if r := compareUint64(v.Minor, v2.Minor); r != 0 {
		return r
	}
	if r := compareUint64(v.Patch, v2.Patch); r != 0 {
		return r
	}
	switch {
	case v.Prerelease == "" && v2.Prerelease == "":
		return 0
	case v.Prerelease == "" && v2.Prerelease != "":
		return 1
	case v.Prerelease != "" && v2.Prerelease == "":
		return -1
	}
	// both prerelease are not equal space
	parts1 := strings.Split(v.Prerelease, ".")
	parts2 := strings.Split(v2.Prerelease, ".")
	l1 := len(parts1)
	l2 := len(parts2)
	l := minInt(l1, l2)

	for i := 0; i < l; i++ {
		p1 := parts1[i]
		p2 := parts2[i]
		u1, err1 := strconv.ParseUint(p1, 10, 64)
		u2, err2 := strconv.ParseUint(p2, 10, 64)
		switch {
		case err1 == nil && err2 == nil: // both are numbers
			if r := compareUint64(u1, u2); r != 0 {
				return r
			}
		case err1 != nil && err2 == nil: // only u2 is numeric
			return 1
		case err1 == nil && err2 != nil: // only u1 is numeric
			return -1

		// both are non numeric
		case p1 > p2:
			return 1
		case p1 < p2:
			return -1
		}
	}

	// all prerelease identifiers are equal -> decide precedence by number of identifiers
	switch {
	case l1 == l2:
		return 0
	case l1 > l2:
		return 1
	default:
		return -1
	}
}
