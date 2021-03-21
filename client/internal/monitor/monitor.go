// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package monitor

import (
	"bytes"
	"time"

	"github.com/stfnmllr/go-resp3/client/internal/conv"
)

// Monitor notification is currently in Redis 6 beta no push type but standard string.

// Notification represents the type for an out of bound monitor push notification send by Redis.
type Notification struct {
	Time time.Time
	Db   int64
	Addr string
	Cmds []string
}

// dirty hack - check if string is monitor push notification
const minNotificationLen = 34

// Parse parses and checks a monitor push notification.
func Parse(b []byte) (*Notification, bool) {
	if len(b) < minNotificationLen {
		return nil, false
	}
	b, seconds, ok := parseInt64(b, []byte{'.'})
	if !ok {
		return nil, false
	}
	b, fraction, ok := parseInt64(b, []byte{' ', '['})
	if !ok {
		return nil, false
	}
	b, db, ok := parseInt64(b, []byte{' '})
	if !ok {
		return nil, false
	}
	b, addr, ok := parseString(b, []byte{']', ' '})
	if !ok {
		return nil, false
	}
	cmds, ok := parseCmd(b)
	if !ok {
		return nil, false
	}
	return &Notification{Time: time.Unix(seconds, nanoSeconds(fraction)), Db: db, Addr: addr, Cmds: cmds}, true
}

func parseInt64(b, sep []byte) ([]byte, int64, bool) {
	pos := bytes.Index(b, sep)
	if pos == -1 {
		return nil, 0, false
	}
	i, err := conv.ParseInt(b[:pos])
	if err != nil {
		return nil, 0, false
	}
	return b[pos+len(sep):], i, true
}

func parseString(b, sep []byte) ([]byte, string, bool) {
	pos := bytes.Index(b, sep)
	if pos == -1 {
		return nil, "", false
	}
	return b[pos+len(sep):], string(b[:pos]), true
}

const (
	pow10_8 = 100000000
	pow10_9 = 1000000000
)

func nanoSeconds(fraction int64) int64 {
	if fraction < pow10_8 {
		for fraction < pow10_8 {
			fraction *= 10
		}
	} else {
		for fraction > pow10_9 {
			fraction /= 10
		}
	}
	return fraction
}

const (
	quote = '"'
	esc   = '\\'
)

func parseQuotedString(b []byte) ([]byte, string, bool) {
	if b[0] != quote {
		return nil, "", false
	}
	buf := make([]byte, len(b)-2)
	j, escaped := 0, false
	for i := 1; i < len(b); i++ {
		switch b[i] {
		case esc:
			escaped = true
		case quote:
			if !escaped {
				return b[i+1:], string(buf[:j]), true
			}
			fallthrough
		default:
			buf[j] = b[i]
			j++
			escaped = false
		}
	}
	return nil, "", false
}

func parseCmd(b []byte) ([]string, bool) {
	r := make([]string, 0, 2)
	var cmd string
	var ok bool
	for {
		b, cmd, ok = parseQuotedString(b)
		if !ok {
			return nil, false
		}
		r = append(r, cmd)
		if len(b) == 0 {
			return r, true
		}
		b = b[1:] // skip blank between commands
	}
}
