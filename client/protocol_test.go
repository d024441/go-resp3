// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"bytes"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestProtocol(t *testing.T) {
	var writeTest = []struct {
		cmd []interface{}
		enc []byte
	}{
		{[]interface{}{1, 2, 3}, []byte("*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n")},
		{[]interface{}{Int64Ptr(1), Int64Ptr(2), Int64Ptr(3)}, []byte("*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n")},
	}

	b := new(bytes.Buffer)
	enc := NewEncoder(b)

	for i, test := range writeTest {
		if err := enc.Encode(test.cmd); err != nil {
			t.Fatalf("line: %d err: %v", i, err)
		}
		if err := enc.Flush(); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(b.Bytes(), test.enc) != 0 {
			t.Fatalf("line: %d got: %v expected: %v", i, b.Bytes(), test.enc)
		}
		b.Reset()
	}
}

func TestDecode(t *testing.T) {
	const bigStr = "3492890328409238509324850943850943825024385"

	bigInt, ok := new(big.Int).SetString(bigStr, 10)
	if !ok {
		t.Fatal("cannot set bigInt")
	}

	var decodeTest = []struct {
		enc []byte
		dec interface{}
	}{
		{[]byte("$0\r\n\r\n"), _string("")},                                        // zero length bulk string
		{[]byte("+Hello World\r\n"), _string("Hello World")},                       // simple string
		{[]byte("$11\r\nHello World\r\n"), _string("Hello World")},                 // bulk string
		{[]byte(":1234\r\n"), _number(1234)},                                       // number
		{[]byte("_\r\n"), _null{}},                                                 // null
		{[]byte(",1.23\r\n"), _double(1.23)},                                       // double
		{[]byte(",10\r\n"), _double(10)},                                           // double as integer
		{[]byte(",inf\r\n"), _double(math.Inf(0))},                                 // double infinite
		{[]byte(",-inf\r\n"), _double(math.Inf(-1))},                               // double neg infinite
		{[]byte("#t\r\n"), _boolean(true)},                                         // boolean true
		{[]byte("#f\r\n"), _boolean(false)},                                        // boolean false
		{[]byte("=15\r\ntxt:Some string\r\n"), _verbatimString("txt:Some string")}, // verbatim string
		{[]byte("(" + bigStr + "\r\n"), (*_bignumber)(bigInt)},                     // big number

		{[]byte("-ERR this is the error description\r\n"), &RedisError{Code: "ERR", Msg: "this is the error description"}}, // simple error
		{[]byte("!21\r\nSYNTAX invalid syntax\r\n"), &RedisError{Code: "SYNTAX", Msg: "invalid syntax"}},                   // blob error

		{[]byte("$?\r\n;4\r\nHell\r\n;6\r\no worl\r\n;1\r\nd\r\n;0\r\n"), _string("Hello world")}, // streamed string
		{ // slice
			[]byte("*4\r\n+first\r\n:1\r\n+second\r\n:2\r\n"),
			_slice{_string("first"), _number(1), _string("second"), _number(2)},
		},
		{ // streamed slice
			[]byte("*?\r\n+a\r\n:1\r\n+b\r\n:2\r\n.\r\n"),
			_slice{_string("a"), _number(1), _string("b"), _number(2)},
		},
		{ // map
			[]byte("%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n"),
			_map{{_string("first"), _number(1)}, {_string("second"), _number(2)}},
		},
		{ // streamed map
			[]byte("%?\r\n+a\r\n:1\r\n+b\r\n:2\r\n.\r\n"),
			_map{{_string("a"), _number(1)}, {_string("b"), _number(2)}},
		},
		{ // set
			[]byte("~2\r\n+first\r\n+second\r\n"),
			_set{_string("first"), _string("second")},
		},
		{ // streamed set
			[]byte("~?\r\n+a\r\n+b\r\n.\r\n"),
			_set{_string("a"), _string("b")},
		},
		{ // attributed slice
			[]byte("|1\r\n+key-popularity\r\n%2\r\n$1\r\na\r\n,0.1923\r\n$1\r\nb\r\n,0.0012\r\n*2\r\n:2039123\r\n:9543892\r\n"),
			attrRedisValue{
				RedisValue: _slice{_number(2039123), _number(9543892)},
				attr:       _map{{_string("key-popularity"), _map{{_string("a"), _double(0.1923)}, {_string("b"), _double(0.0012)}}}},
			},
		},
		{ // attributed slice element
			[]byte("*3\r\n:1\r\n:2\r\n|1\r\n+ttl\r\n:3600\r\n:3\r\n"),
			_slice{
				_number(1),
				_number(2),
				attrRedisValue{
					RedisValue: _number(3),
					attr:       _map{{_string("ttl"), _number(3600)}},
				},
			},
		},
	}

	b := new(bytes.Buffer)
	dec := NewDecoder(bufio.NewReader(b))

	for i, test := range decodeTest {
		b.Reset()
		if _, err := b.Write(test.enc); err != nil {
			t.Fatalf("line: %d err: %v", i, err)
		}

		dec, err := dec.Decode()
		if err != nil {
			t.Fatalf("line: %d err: %v", i, err)
		}

		if !reflect.DeepEqual(dec, test.dec) {
			t.Fatalf("line: %d value got %#v - expected %#v", i, dec, test.dec)
		}
	}
}
