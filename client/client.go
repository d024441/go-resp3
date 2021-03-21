// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net"
	"os"
)

//go:generate converter

// Client related attributes like name and version.
const (
	ClientName    = "go-resp3"
	ClientVersion = "0.11.5"
)

// Environment variables.
const (
	EnvHost = "REDIS_HOST"
	EnvPort = "REDIS_PORT"
)

// Default redis host and port.
const (
	LocalHost        = "127.0.0.1"
	DefaultRedisPort = "6379"
)

const (
	tcpNetwork = "tcp"
)

func hostPort(address string) string {
	if address != "" { // use provided address
		return address
	}
	host, ok := os.LookupEnv(EnvHost)
	if !ok {
		host = LocalHost
	}
	port, ok := os.LookupEnv(EnvPort)
	if !ok {
		port = DefaultRedisPort
	}
	return net.JoinHostPort(host, port)
}
