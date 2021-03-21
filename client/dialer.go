// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"time"
)

// InvalidateCallback is the function type for the Redis cache invalidate callback function.
type InvalidateCallback func(keys []string)

// MonitorCallback is the function type for the Redis monitor callback function.
type MonitorCallback func(time time.Time, db int64, addr string, cmds []string)

// TraceCallback is the function type for the tracing callback function.
type TraceCallback func(dir bool, b []byte)

// SendInterceptor is the function type for the send interceptor function.
type SendInterceptor func(name string, values []interface{})

//Dialer contains options for connecting to a redis server.
type Dialer struct {
	net.Dialer
	// TLS
	TLSConfig *tls.Config
	// Connection logging.
	Logger *log.Logger
	// Channel size for asynchronous result reader and handler.
	ChannelSize int
	// Duration to wait for async results before timeout.
	AsyncTimeout time.Duration
	// Redis authentication.
	Username, Password string
	// Redis client name.
	ClientName string
	// Client cache invalidation callback.
	InvalidateCallback InvalidateCallback
	// Monitor callback.
	MonitorCallback MonitorCallback
	// Callback function tracing Redis commands and results on RESP3 protocol level.
	// Direction dir is true for sent bytes b (commands), false for received bytes b (results).
	TraceCallback TraceCallback
	// Command interceptor (debugging).
	SendInterceptor SendInterceptor
}

func (d *Dialer) channelSize() int {
	if d.ChannelSize < 1 {
		return defaultChannelSize
	}
	if d.ChannelSize < minimumChannelSize {
		return minimumChannelSize
	}
	return d.ChannelSize
}

// Dial connects to the redis server address.
func (d *Dialer) Dial(address string) (Conn, error) {
	return d.DialContext(context.Background(), address)
}

// DialContext connects to the redis server address using the provided context.
func (d *Dialer) DialContext(ctx context.Context, address string) (Conn, error) {
	return d.dialContext(ctx, address)
}

func (d *Dialer) dialContext(ctx context.Context, address string) (*conn, error) {
	c, err := d.Dialer.DialContext(ctx, tcpNetwork, hostPort(address))
	if err != nil {
		return nil, err
	}
	if d.TLSConfig != nil {
		c = tls.Client(c, d.TLSConfig)
	}
	return newConn(nil, c, d)
}
