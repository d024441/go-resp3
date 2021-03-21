// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

// This example shows how to create custom connection and pipeline objects
// to extend the command interface or redefine commands.

import (
	"fmt"
	"log"

	"github.com/stfnmllr/go-resp3/client"
)

// MyCommands defines redefined and additional commands.
type MyCommands interface {
	// Redefinition of Del - use variadic attribute instead slice.
	Del(keys ...interface{}) client.Result
	// Additional command Heartbeat.
	Heartbeat() client.Result
}

// Check, that connection and pipeline implement MyCommands.
var (
	_ MyCommands = (*MyConn)(nil)
	_ MyCommands = (*MyPipeline)(nil)
)

// MyConn is the custom connection.
type MyConn struct {
	client.Conn
}

// MyConn contructor.
func NewMyConn() (*MyConn, error) {
	conn, err := client.Dial("")
	if err != nil {
		return nil, err
	}
	return &MyConn{Conn: conn}, nil
}

// Pipeline overwrites the connection pipeline method.
func (c *MyConn) Pipeline() *MyPipeline {
	return NewMyPipeline(c)
}

// Del implements MyCommands interface.
func (c *MyConn) Del(keys ...interface{}) client.Result {
	// call client connection method.
	return c.Conn.Del(keys)
}

// Heartbeat implements MyCommands interface.
func (c *MyConn) Heartbeat() client.Result {
	// call client connection method
	// using the generic Do method
	return c.Conn.Do("ping", "ok")
}

// MyPipeline is the custom pipeline.
type MyPipeline struct {
	client.Pipeline
}

// MyPipeline constructor.
func NewMyPipeline(c *MyConn) *MyPipeline {
	return &MyPipeline{Pipeline: c.Conn.Pipeline()}
}

// Del implements MyCommands interface.
func (p *MyPipeline) Del(keys ...interface{}) client.Result {
	return p.Pipeline.Del(keys)
}

// Heartbeat implements MyCommands interface.
func (p *MyPipeline) Heartbeat() client.Result {
	return p.Pipeline.Do("ping", "ok")
}

func Example_redefine() {
	// Use custom connection.
	conn, err := NewMyConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.Set("mykey", "myvalue").Err(); err != nil {
		log.Fatal(err)
	}
	// Use of redefined connection method.
	removedKeys, err := conn.Del("mykey").ToInt64()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(removedKeys)

	// Use of extended connection method.
	ok, err := conn.Heartbeat().ToString()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ok)

	// Use custom pipeline.
	pipeline := conn.Pipeline()

	set := pipeline.Set("mykey", "myvalue")
	// Use of redefined pipeline method.
	del := pipeline.Del("mykey")
	// Use of extended pipeline method.
	heartbeat := conn.Heartbeat()

	pipeline.Flush()

	if err := set.Err(); err != nil {
		log.Fatal(err)
	}
	removedKeys, err = del.ToInt64()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(removedKeys)

	ok, err = heartbeat.ToString()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ok)
	// Output:
	// 1
	// ok
	// 1
	// ok
}
