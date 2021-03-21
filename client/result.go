// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"
)

// ErrNotFlushed is returned by a command executed in a unflushed pipeline.
var ErrNotFlushed = errors.New("pipeline: not flushed")

// ErrTimeout is returned by a command after timeout is reached before a async Redis result is available.
var ErrTimeout = errors.New("timeout while waiting for result")

// Result represents a redis command result.
type Result interface {
	// Attr returns the attribute of a Redis value if provided - <nil> otherwise.
	Attr() (*Map, error)
	// Err returns any of the following errors:
	// - <nil>: result received successfully.
	// - ErrNotFlushed: result not received yet (pipeline not flushed).
	// - ErrTimeout: timeout reached before result is available.
	// - RedisError: redis error message if a redis command was executed unsuccessfully.
	Err() error
	// IsNull returns <true> if the redis value is null.
	IsNull() (bool, error)
	// Kind returns the type of a Redis value.
	Kind() (RedisKind, error)
	// Value returns a Redis value.
	Value() (RedisValue, error)
	// Conversion methods.
	Converter
}

const (
	rsNotFlushed uint32 = 0
	rsFlushed    uint32 = 1 << iota
	rsAvailable
	rsWaiting
	rsSetting
)

type result struct {
	// caution: field alignment (memory consumption)
	value   RedisValue
	err     error
	request *request
	flags   uint32
}

func newResult() *result {
	return &result{request: freeRequest.get()}
}

func (r *result) cmd() []interface{} {
	return r.request.cmd
}

// set error before flush like InvalidValueError
func (r *result) setErr(err error) {
	// check on not flushed
	if atomic.LoadUint32(&r.flags) != rsNotFlushed {
		panic("cannot set error for result in not flushed state")
	}
	r.err = err
	atomic.StoreUint32(&r.flags, rsAvailable)
}

func (r *result) flush() {
	atomic.StoreUint32(&r.flags, rsFlushed)
}

func (r *result) wait() error {
	if atomic.LoadUint32(&r.flags) == rsNotFlushed {
		return ErrNotFlushed
	}

	for !atomic.CompareAndSwapUint32(&r.flags, rsFlushed, rsWaiting) {
		if atomic.LoadUint32(&r.flags) == rsAvailable {
			return r.err
		}
		runtime.Gosched() // spinning - allow other goroutines to process
	}

	// wait for done signal
	<-r.request.done
	if atomic.LoadUint32(&r.flags) != rsAvailable {
		panic("result - inconsistent state")
	}
	return r.err
}

func (r *result) ack(value RedisValue, err error) {
	// todo - state type + check state

	isWaiting := !atomic.CompareAndSwapUint32(&r.flags, rsFlushed, rsSetting)
	r.value = value
	r.err = err
	atomic.StoreUint32(&r.flags, rsAvailable)
	if isWaiting {
		r.request.done <- true
	}
	// free data
	freeRequest.put(r.request)
	r.request = nil
}

func (r *result) setTimeout(timeout time.Duration) {
	r.request.timeout = timeout
}

// Err returns any of the following errors:
// - <nil>: result received successfully.
// - ErrNotFlushed: result not received yet (pipeline not flushed).
// - ErrTimeout: timeout reached before result is available.
// - RedisError: redis error message if a redis command was executed unsuccessfully.
func (r *result) Err() error {
	if err := r.wait(); err != nil {
		return err
	}
	return r.err
}

// Kind returns the type of a Redis value.
func (r *result) Kind() (RedisKind, error) {
	if err := r.wait(); err != nil {
		return RkInvalid, err
	}
	return r.value.Kind(), nil
}

// Value returns a Redis value.
func (r *result) Value() (RedisValue, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	return r.value, nil
}

// Attr returns the attribute of a Redis value if provided - <nil> otherwise.
func (r *result) Attr() (*Map, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	return r.value.Attr(), nil
}

// IsNull returns <true> if the redis value is null.
func (r *result) IsNull() (bool, error) {
	if err := r.wait(); err != nil {
		return false, err
	}
	return r.value.Kind() == RkNull, nil
}

// test on Result interface implementation
var _ Result = (*result)(nil)
