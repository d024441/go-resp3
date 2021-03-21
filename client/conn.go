// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stfnmllr/go-resp3/client/internal/monitor"
)

// MsgCallback is the function type for Redis pubsub message callback functions.
type MsgCallback func(pattern, channel, msg string)

type sendFct func(name string, r *result)

// ErrInShutdown is returned by a command executed during client is in shutdown status.
var ErrInShutdown = errors.New("connection is shutdown")

// Conn represents the redis network connection.
type Conn interface {
	Commands
	Pipeline() Pipeline
	Close() error
	ConnInfo() ConnInfo
	private() // private interface
}

// ConnInfo provided information about a connection.
type ConnInfo struct {
	RedisVersion Version
}

// ErrConnClosed is returned when calling methods on connection after the connection is closed.
var ErrConnClosed = errors.New(ClientName + ": database is already closed")

// Dial connects to the redis server address.
func Dial(address string) (Conn, error) {
	return new(Dialer).Dial(address)
}

// DialContext connects to the redis server address using the provided context.
func DialContext(ctx context.Context, address string) (Conn, error) {
	return new(Dialer).DialContext(ctx, address)
}

// check interface implementations.
var (
	_ Conn = (*conn)(nil)
)

type conn struct {
	inShutdown int32 // atomic access

	mu       sync.RWMutex
	encodemu sync.Mutex // protects encoding + flush

	db *db // connection owned by db - nil otherwise

	netConn net.Conn
	logger  *log.Logger

	closed, pooled bool

	*command

	dec Decoder
	enc Encoder

	readChan chan interface{}
	resChan  chan []*result
	sendChan chan *result

	hello Result

	asyncTimeout       time.Duration
	invalidateCallback InvalidateCallback
	monitorCallback    MonitorCallback

	sendInterceptor SendInterceptor

	nextResult func() *result

	shutdown <-chan bool
}

const (
	defaultChannelSize = 10000
	minimumChannelSize = 100
)

func newConn(db *db, netConn net.Conn, d *Dialer) (*conn, error) {
	c := &conn{
		db:                 db,
		netConn:            netConn,
		logger:             d.Logger,
		readChan:           make(chan interface{}, d.channelSize()),
		resChan:            make(chan []*result, defResultItems),
		sendChan:           make(chan *result, d.channelSize()),
		asyncTimeout:       d.AsyncTimeout,
		invalidateCallback: d.InvalidateCallback,
		monitorCallback:    d.MonitorCallback,
		sendInterceptor:    d.SendInterceptor,
	}

	if c.logger != nil {
		c.logger.Printf("remote address %s - local address %s", c.netConn.RemoteAddr().String(), c.netConn.LocalAddr().String())
	}

	c.command = newCommand(c.send, c.sendInterceptor)

	if d.TraceCallback != nil {
		c.enc, c.dec = tracer(d.TraceCallback, netConn)
	} else {
		c.enc = NewEncoder(netConn)
		c.dec = NewDecoder(netConn)
	}

	c.nextResult = c.resultIterator()

	c.shutdown = c.watch()

	var auth *UsernamePassword
	if d.Username != "" || d.Password != "" {
		auth = &UsernamePassword{Username: d.Username, Password: d.Password}
	}
	var clientName *string
	if d.ClientName != "" {
		clientName = &d.ClientName
	}

	c.hello = c.Hello(protocolVersion, auth, clientName)
	if err := c.hello.Err(); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (c *conn) private() {} // private interface

func (c *conn) resultIterator() func() *result {
	var (
		results   []*result
		size, pos int
	)

	return func() *result {
		if pos >= size {
			if results != nil {
				freeResults.put(results)
			}
			pos = 0
			var ok bool
			for {
				results, ok = <-c.resChan
				if !ok {
					return nil
				}
				// skip empty list
				if size = len(results); size != 0 {
					break
				}
			}
		}
		i := pos
		pos++
		return results[i]
	}
}

func (c *conn) send(name string, r *result) {
	if atomic.LoadInt32(&c.inShutdown) != 0 {
		r.setErr(ErrInShutdown)
		return
	}
	r.flush() // no pipeline
	c.sendChan <- r
}

const (
	helpVersion = "version"
)

func (c *conn) ConnInfo() ConnInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ci := ConnInfo{}

	if m, err := c.hello.ToStringMap(); err == nil {
		if s, ok := m[helpVersion].(string); ok {
			ci.RedisVersion = ParseVersion(s)
		}
	}
	return ci
}

func (c *conn) quitLocked() error {
	return c.Quit().Err()
}

func (c *conn) waitClosedLocked() error {
	c.closed = true
	select {
	case <-c.shutdown: // wait for watcher to shutdown
		return c.netConn.Close()
		// TODO timeout
	}
}

func (c *conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed || c.pooled { // pooled - second attempt to close
		return ErrConnClosed
	}

	if c.db != nil && c.db.putConn(c) { // connection owned by db
		c.pooled = true
		return nil
	}
	c.quitLocked()
	return c.waitClosedLocked()
}

func (c *conn) Pipeline() Pipeline {
	return newPipeline(c)
}

func (c *conn) watch() <-chan bool {

	shutdown := make(chan bool)

	go func() {

		var wgReader sync.WaitGroup  // wait for reader
		var wgHandler sync.WaitGroup // wait for handler
		var wgSender sync.WaitGroup  // wait for reader

		var errChan = make(chan error)

		wgHandler.Add(1)
		go c.cmdHandler(&wgHandler, c.readChan)
		wgReader.Add(1)
		go c.reader(&wgReader, c.readChan, errChan)
		wgSender.Add(1)
		go c.sender(&wgSender, c.sendChan)

		err := <-errChan
		wgReader.Wait() // wait for reader

		atomic.StoreInt32(&c.inShutdown, 1)

		close(c.readChan) // stop handler
		wgHandler.Wait()  // wait for handler

		close(c.sendChan) // stop sender
		wgSender.Wait()   // wait for sender

		// drain resChan
		close(c.resChan)
		for {
			r := c.nextResult()
			if r == nil {
				break
			}
			r.ack(nil, err)
		}

		close(shutdown)

	}()

	return shutdown

}

func (c *conn) cmdHandler(wg *sync.WaitGroup, readChan <-chan interface{}) {
	channelMap := map[string]MsgCallback{}

	defer wg.Done()

	for {

		val, ok := <-readChan
		if !ok { // channel closed
			return
		}

		switch val := val.(type) {

		case RedisValue:
			result := c.nextResult()
			result.ack(val, nil)

		case error:
			result := c.nextResult()
			result.ack(nil, val)

		case *subscribeNotification:
			if ok := c.handleSubscribeNotification(val, channelMap, readChan); !ok {
				return
			}

		case *unsubscribeNotification:
			if ok := c.handleUnsubscribeNotification(val, channelMap, readChan); !ok {
				return
			}

		case *publishNotification:
			if cb, ok := channelMap[val.channel]; ok && cb != nil {
				cb(val.pattern, val.channel, val.msg)
			}

		case *invalidateNotification:
			if c.invalidateCallback != nil {
				c.invalidateCallback(val.keys)
			}

		case *monitor.Notification:
			if c.monitorCallback != nil {
				c.monitorCallback(val.Time, val.Db, val.Addr, val.Cmds)
			}

		case *genericNotification:

		default:
			panic("invalid messsage type")
		}
	}
}

func (c *conn) handleSubscribeNotification(n *subscribeNotification, channelMap map[string]MsgCallback, readChan <-chan interface{}) bool {
	result := c.nextResult()

	size := len(result.request.cmd)
	channels := result.request.cmd[1:size]

	for i, ch := range channels { // expect a push message for all subscribed channels

		ch := ch.(string)

		if i != 0 {
			val, ok := <-readChan
			if !ok { // channel closed
				return false
			}
			n = val.(*subscribeNotification)
		}

		if n.channel != ch {
			panic("subscribe: command message channel mismatch")
		}
		channelMap[ch] = result.request.cb
	}

	result.ack(nil, nil)
	return true
}

func (c *conn) handleUnsubscribeNotification(n *unsubscribeNotification, channelMap map[string]MsgCallback, readChan <-chan interface{}) bool {
	result := c.nextResult()

	size := len(result.request.cmd)
	channels := result.request.cmd[1:size]

	if size > 1 { // unsubscribe list of channels

		for i, ch := range channels { // expect a push message for all unsubscribed channels

			ch := ch.(string)

			if i != 0 {
				val, ok := <-readChan
				if !ok { // channel closed
					return false
				}
				n = val.(*unsubscribeNotification)
			}
			if n.channel != ch {
				panic("unsubscribe: command message channel mismatch")
			}
			delete(channelMap, ch)
		}

		result.ack(nil, nil)
		return true
	}

	// unsubscribe from all channels
	for {
		delete(channelMap, n.channel)
		if n.count == 0 {
			break
		}
		val, ok := <-readChan
		if !ok { // channel closed
			return false
		}
		n = val.(*unsubscribeNotification)
	}

	result.ack(nil, nil)
	return true
}

func (c *conn) reader(wg *sync.WaitGroup, readChan chan<- interface{}, errorChan chan<- error) {
	defer wg.Done()

	for {
		val, err := c.dec.Decode()
		if err != nil {
			errorChan <- err
			return
		}
		readChan <- val
	}
}

func (c *conn) flush(pipeline bool, results []*result) error {
	c.encodemu.Lock()
	defer c.encodemu.Unlock()
	for _, r := range results {
		if pipeline {
			r.flush()
		}
		r.setTimeout(c.asyncTimeout)
		c.enc.Encode(r.cmd())
	}
	if err := c.enc.Flush(); err != nil {
		panic(err) // TODO
	}
	c.resChan <- results
	return nil
}

func (c *conn) sender(wg *sync.WaitGroup, sendChan <-chan *result) {
	defer wg.Done()

	results := freeResults.get()

	for {
		select {
		case r, ok := <-sendChan:
			if !ok {
				goto close
			}
			results = append(results, r)

			burst := true
			for burst {
				select {
				case r, ok := <-sendChan:
					if !ok {
						goto close
					}
					results = append(results, r)
				default:
					burst = false
				}
			}
			c.flush(false, results)
			results = freeResults.get()
		}
	}
close:
	c.flush(false, results)
}
