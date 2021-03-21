// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stfnmllr/go-resp3/client"
)

func lpadName(name string, cnt, n int) string {
	num := strconv.Itoa(cnt)
	return name + strings.Repeat("_", n-len(num)) + num
}

func benchmarkSetGet(conn client.Conn, cnt int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < cnt; j++ {
			if err := conn.Set("foo", "bar").Err(); err != nil {
				b.Fatal(err)
			}
			value, err := conn.Get("foo").ToString()
			if err != nil {
				b.Fatal(err)
			}
			if value != "bar" {
				b.Fatal("got wrong value")
			}
		}
	}
}

func benchmarkSetGetAsync(conn client.Conn, cnt int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		results := make([]client.Result, 2*cnt)
		for j := 0; j < cnt; j++ {
			results[j*2] = conn.Set("foo", "bar")
			results[j*2+1] = conn.Get("foo")
		}
		for j := 0; j < cnt; j++ {
			if err := results[j*2].Err(); err != nil {
				b.Fatal(err)
			}
			value, err := results[j*2+1].ToString()
			if err != nil {
				b.Fatal(err)
			}
			if value != "bar" {
				b.Fatal("got wrong value")
			}
		}
	}
}

func benchmarkPipelineSetGet(conn client.Conn, cnt int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		results := make([]client.Result, 2*cnt)
		p := conn.Pipeline()
		for j := 0; j < cnt; j++ {
			results[j*2] = p.Set("foo", "bar")
			results[j*2+1] = p.Get("foo")
		}
		if err := p.Flush(); err != nil {
			b.Fatal(err)
		}
		for j := 0; j < cnt; j++ {
			if err := results[j*2].Err(); err != nil {
				b.Fatal(err)
			}
			value, err := results[j*2+1].ToString()
			if err != nil {
				b.Fatal(err)
			}
			if value != "bar" {
				b.Fatal("got wrong value")
			}
		}
	}
}

func BenchmarkCommand(b *testing.B) {
	dialer := client.Dialer{Logger: log.New(os.Stderr, "", log.LstdFlags)}
	conn, err := dialer.Dial("")
	if err != nil {
		b.Fatal(err)
	}

	b.Run("SingleCall", func(b *testing.B) {
		const x = 6 // 10^x
		cnt := 1
		for i := 0; i < x; i++ {
			b.Run(lpadName("SetGet", cnt, x), func(b *testing.B) { benchmarkSetGet(conn, cnt, b) })
			cnt *= 10
		}
	})

	b.Run("SingleCallAsync", func(b *testing.B) {
		const x = 6 // 10^x
		cnt := 1
		for i := 0; i < x; i++ {
			b.Run(lpadName("SetGet", cnt, x), func(b *testing.B) { benchmarkSetGetAsync(conn, cnt, b) })
			cnt *= 10
		}
	})

	b.Run("Pipeline", func(b *testing.B) {
		const x = 7 // 10^x
		cnt := 1
		for i := 0; i < x; i++ {
			b.Run(lpadName("SetGet", cnt, x), func(b *testing.B) { benchmarkPipelineSetGet(conn, cnt, b) })
			cnt *= 10
		}
	})

	if err := conn.Close(); err != nil {
		b.Fatal(err)
	}
}
