// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultMaxIdleConns  = 2
	closeConnChannelSize = 10
)

// ErrDBClosed is returned when calling methods on database after the database is closed.
var ErrDBClosed = errors.New(ClientName + ": database is already closed")

// DB is a database holding a pool of redis connections.
// *** not yet completely implemented - experimental ***
type DB interface {
	Commands
	Conn(ctx context.Context) (Conn, error)
	// Pipeline() Pipeline
	Close() error
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
	private() // private interface
}

// OpenDB open a new database.
func OpenDB(address string, dialer *Dialer) DB {
	return newDB(address, dialer)
}

// check interface implementations.
var (
	_ DB = (*db)(nil)
)

type db struct {
	// stats (please see https://golang.org/pkg/database/sql for reference)
	// atomic access (64-bit alignment on 32 platforms - please see https://golang.org/pkg/sync/atomic)
	waitDuration int64 // Wait duration.

	address string
	dialer  *Dialer
	closed  int32

	// connection pooling attributes
	mu          sync.RWMutex
	freeConn    []*conn
	waitConn    map[chan *conn]struct{}
	closeConnCh chan []*conn
	wg          sync.WaitGroup

	// stats (please see https://golang.org/pkg/database/sql for reference)
	numOpen int // number of opened and pending open connections
	maxIdle int // zero means defaultMaxIdleConns; negative means 0
	maxOpen int // <= 0 means unlimited
	//maxLifetime       time.Duration // maximum amount of time a connection may be reused
	waitCount         int64 // Total number of connections waited for.
	maxIdleClosed     int64 // Total number of connections closed due to idle.
	maxLifetimeClosed int64 // Total number of connections closed due to max free limit.

	*command
}

func newDB(address string, dialer *Dialer) *db {
	if dialer == nil {
		dialer = new(Dialer)
	}
	db := &db{
		address:     address,
		dialer:      dialer,
		freeConn:    make([]*conn, 0),
		waitConn:    make(map[chan *conn]struct{}),
		closeConnCh: make(chan []*conn, closeConnChannelSize),
		maxIdle:     defaultMaxIdleConns,
	}
	db.command = newCommand(db.send, nil)

	db.wg.Add(1)
	go db.connCloser(&db.wg, db.closeConnCh)

	return db
}

func (db *db) private() {} // private interface

func (db *db) Close() error { return db.close() }

// TODO check inflight

func (db *db) Conn(ctx context.Context) (Conn, error) { return db.getConn(ctx) }

func (db *db) SetConnMaxLifetime(d time.Duration) { panic("not implemented") }

func (db *db) SetMaxIdleConns(n int) {
	if n < 0 {
		n = 0
	}
	db.mu.Lock()
	db.maxIdle = n
	if db.maxOpen > 0 && db.maxIdle > db.maxOpen {
		db.maxIdle = db.maxOpen
	}
	if len(db.freeConn) <= db.maxIdle {
		db.mu.Unlock()
		return
	}

	// shrink freeConn and close obsolete connections
	closeConn := db.freeConn[db.maxIdle:]
	db.freeConn = db.freeConn[:db.maxIdle]

	db.maxIdleClosed += int64(len(closeConn))
	db.mu.Unlock()

	db.closeConnCh <- closeConn
}

func (db *db) SetMaxOpenConns(n int) {
	if n < 0 {
		n = 0
	}
	db.mu.Lock()
	db.maxOpen = n
	correctMaxIdle := db.maxOpen > 0 && db.maxIdle > db.maxOpen
	db.mu.Unlock()
	if correctMaxIdle {
		db.SetMaxIdleConns(n)
	}
}

func (db *db) Stats() sql.DBStats { return db.dbStats() }

func (db *db) send(name string, r *result) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO timeout configuration
	defer cancel()

	conn, err := db.getConn(ctx)
	if err != nil {
		r.setErr(err)
		return
	}
	defer db.putConn(conn)
	conn.send(name, r)
}

func (db *db) dbStats() sql.DBStats {
	wait := atomic.LoadInt64(&db.waitDuration)

	db.mu.RLock()
	defer db.mu.RUnlock()
	numFree := len(db.freeConn)
	return sql.DBStats{
		MaxOpenConnections: db.maxOpen,

		Idle:            numFree,
		OpenConnections: db.numOpen,
		InUse:           db.numOpen - numFree,

		WaitCount:         db.waitCount,
		WaitDuration:      time.Duration(wait),
		MaxIdleClosed:     db.maxIdleClosed,
		MaxLifetimeClosed: db.maxLifetimeClosed,
	}
}

func (db *db) getConn(ctx context.Context) (*conn, error) {
	db.mu.Lock()
	numFree := len(db.freeConn)

	switch {

	case (db.maxOpen <= 0 || db.numOpen <= db.maxOpen) && numFree > 0:
		conn := db.freeConn[0]
		copy(db.freeConn, db.freeConn[1:]) //fifo
		db.freeConn = db.freeConn[:numFree-1]
		db.mu.Unlock()
		return conn, nil

	case db.numOpen <= 0 || db.numOpen <= db.maxOpen:
		conn, err := db.dialer.dialContext(ctx, db.address)
		if err != nil {
			db.mu.Unlock()
			return nil, err
		}
		db.numOpen++
		db.mu.Unlock()
		return conn, nil

	}

	// wait for connection
	db.waitCount++
	r := make(chan *conn, 1)
	db.waitConn[r] = struct{}{}
	db.mu.Unlock()

	waitStart := time.Now()

	select {
	case <-ctx.Done(): // context

		atomic.AddInt64(&db.waitDuration, int64(time.Since(waitStart)))

		db.mu.Lock()
		delete(db.waitConn, r) // delete request
		db.mu.Unlock()

		select {
		default:
		case conn := <-r: // check if connection was already provided
			db.putConn(conn)
		}
		db.mu.Unlock()

		return nil, ctx.Err()
	case conn := <-r: // connection
		atomic.AddInt64(&db.waitDuration, int64(time.Since(waitStart)))
		return conn, nil
	}
}

func (db *db) putConn(c *conn) bool {
	if c == nil {
		panic("putConn: connection is nil") // should never happen
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// check if somebody is waiting for a connection
	for r := range db.waitConn {
		delete(db.waitConn, r)
		r <- c
		return true
	}

	if db.maxOpen > 0 && db.numOpen >= db.maxOpen {
		db.numOpen--
		return false
	}
	if len(db.freeConn) >= db.maxIdle {
		db.numOpen--
		return false
	}
	db.freeConn = append(db.freeConn, c)
	return true
}

func (db *db) close() error {
	if !atomic.CompareAndSwapInt32(&db.closed, 0, 1) {
		return ErrDBClosed
	}

	// close idle connections
	db.closeConnCh <- db.freeConn

	// stop connCloser
	close(db.closeConnCh)
	// wait for go routines to stop
	db.wg.Wait()
	return nil
}

func (db *db) connCloser(wg *sync.WaitGroup, closeConnCh <-chan []*conn) {
	for closeConn := range closeConnCh {
		for _, c := range closeConn {
			c.mu.Lock()
			if !c.closed {
				c.quitLocked()
			}
		}
		for _, c := range closeConn {
			c.waitClosedLocked()
			c.mu.Unlock()
		}
	}
	wg.Done()
}
