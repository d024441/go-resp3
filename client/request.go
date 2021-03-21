// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"sync"
	"time"
)

const (
	minRequests = 2000
	maxRequests = 20000
	defCmdSize  = 5
)

var freeRequest = requestPool{}

func init() {
	for i := 0; i < minRequests; i++ {
		freeRequest.put(newRequest())
	}
}

type requestPool struct {
	mu   sync.Mutex
	size int
	free *request
}

func (p *requestPool) get() *request {
	p.mu.Lock()
	if p.free == nil {
		p.mu.Unlock()
		return newRequest()
	}
	p.size--
	r := p.free
	p.free = r.next
	p.mu.Unlock()
	return r
}

func (p *requestPool) put(r *request) {
	p.mu.Lock()
	if p.size >= maxRequests {
		p.mu.Unlock()
		return
	}
	r.cb = nil
	r.cmd = r.cmd[:0]
	p.size++
	r.next = p.free
	p.free = r
	p.mu.Unlock()
}

type request struct {
	cmd     []interface{} // Redis command 'token'
	done    chan bool
	cb      MsgCallback // pubsub callback function
	timeout time.Duration
	next    *request
}

func newRequest() *request {
	return &request{
		cmd:  make([]interface{}, 0, defCmdSize),
		done: make(chan bool),
	}
}
