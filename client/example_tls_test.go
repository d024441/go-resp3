// +build TLS

// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/stfnmllr/go-resp3/client"
	"github.com/stfnmllr/go-resp3/internal/ioutildepr"
)

/*
This example is based on the TLS Support getting started of Redis 6.
Please see https://github.com/antirez/redis/blob/unstable/TLS.md for details.

The go-resp3 client TLS support is solely based on Go standard library capabilities, especially the
- TLS configuration (https://golang.org/pkg/crypto/tls/#Config) and
- client connection (https://golang.org/pkg/crypto/tls/#Client)

Please check the Go standard library documentation for further information.
*/

const (
	certFile   = "tls/redis.crt"
	keyFile    = "tls/redis.key"
	cacertFile = "tls/ca.crt"
)

func Example_TLS() {
	// Root certificate.
	cacert, err := ioutildepr.ReadFile(cacertFile)
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(cacert) {
		log.Fatal("failed to parse root certificate")
	}

	// Key pair.
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create TLS configuration.
	config := &tls.Config{
		// whether the ServerName needs to be provided or InsecureSkipVerify needs to be set
		// please see the comments in https://golang.org/pkg/crypto/tls/#Client
		InsecureSkipVerify: true,
		RootCAs:            roots,
		Certificates:       []tls.Certificate{cert},
	}

	// Create dialer with TLS configuration.
	dialer := &client.Dialer{TLSConfig: config}

	// Open connection with standard host and port.
	conn, err := dialer.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	if err := conn.Set("mykey", "Hello Redis").Err(); err != nil {
		log.Fatal(err)
	}
	value, err := conn.Get("mykey").ToString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
	// Output:
	// Hello Redis
}
