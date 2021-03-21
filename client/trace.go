// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"io"
	"sync"
)

type trace struct {
	cb  TraceCallback
	b   []byte
	dec Decoder
	mu  sync.Mutex
}

func tracer(cb TraceCallback, rw io.ReadWriter) (Encoder, Decoder) {
	t := &trace{cb: cb, b: make([]byte, 0, 64)}

	encWriter := newTraceWriter(t.writeEnc)
	decWriter := newTraceWriter(t.writeDec)

	enc := NewEncoder(io.MultiWriter(rw, encWriter))

	// !Caution: use 'own' bufio.Writer (even Decoder would create its own).
	// We need to be sure, that the TeeReader is based on bufio.Reader and NOT based on
	// the connection itself as other bufio.Reader (Decoder) might 'read ahead' on connection.
	t.dec = NewDecoder(io.TeeReader(bufio.NewReader(rw), decWriter))

	//return tracer as Decoder.
	return enc, t
}

func (t *trace) call(dir bool, b []byte) {
	t.mu.Lock()
	t.cb(dir, b)
	t.mu.Unlock()
}

func (t *trace) Decode() (interface{}, error) {
	value, err := t.dec.Decode()
	if err == nil {
		t.call(false, t.b)
	}
	t.b = t.b[:0]
	return value, err
}

func (t *trace) writeEnc(p []byte) (int, error) {
	t.call(true, p)
	return len(p), nil
}

func (t *trace) writeDec(p []byte) (int, error) {
	t.b = append(t.b, p...)
	return len(p), nil
}

type traceWriter struct {
	fct func(p []byte) (int, error)
}

func newTraceWriter(fct func(p []byte) (int, error)) *traceWriter {
	return &traceWriter{fct: fct}
}

func (w *traceWriter) Write(p []byte) (int, error) {
	return w.fct(p)
}
