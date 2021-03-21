// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"fmt"
	"log"

	"time"

	"github.com/stfnmllr/go-resp3/client"
)

func Monitor(time time.Time, db int64, addr string, cmd []string) {
	fmt.Printf("time: %s database: %d address: %s command %v\n", time, db, addr, cmd)
}

func Example_monitor() {
	// Register monitor callback.
	dialer := client.Dialer{MonitorCallback: Monitor}
	conn, err := dialer.Dial("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Start monitor.
	ok, err := conn.Monitor().ToBool()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Start monitor failed")
	}

	// Process some commands in same connection.
	mykey := client.RandomKey("mykey")
	for i := 0; i < 5; i++ {
		if err = conn.Set(mykey, "myValue").Err(); err != nil {
			log.Fatal(err)
		}
	}
}
