// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"sync/atomic"
)

// Pipeline supports redis pipelining capabilities.
// Multiple goroutines must not invoke methods on a Pipeline simultaneously.
type Pipeline interface {
	Commands
	Reset()
	Flush() error
}

var _ Pipeline = (*pipeline)(nil)

type pipeline struct {
	c       *conn
	err     error
	results []*result
	*command
}

func newPipeline(c *conn) *pipeline {
	p := &pipeline{c: c, results: freeResults.get()}
	p.command = newCommand(p.send, c.sendInterceptor)
	return p
}

func (p *pipeline) send(name string, r *result) {
	if atomic.LoadInt32(&p.c.inShutdown) != 0 {
		r.setErr(ErrInShutdown)
		return
	}
	p.results = append(p.results, r)
}

func (p *pipeline) Reset() {
	p.results = p.results[:0]
}

func (p *pipeline) Flush() error {
	err := p.c.flush(true, p.results)
	p.results = freeResults.get()
	return err
}
