// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"fmt"
	"log"

	"github.com/stfnmllr/go-resp3/client"
)

func Example_clientcache() {
	mykey := client.RandomKey("mykey")

	invalidated := make(chan struct{}, 0)

	// Create connetion providing key invalidation callback.
	dialer := new(client.Dialer)
	dialer.InvalidateCallback = func(keys []string) {
		for _, key := range keys {
			if key == mykey {
				close(invalidated)
				fmt.Println("Key invalidated")
				break
			}
		}
	}

	conn, err := dialer.Dial("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Set client tracking on.
	if err := conn.ClientTracking(true, nil, nil, false, false, false, false).Err(); err != nil {
		log.Fatal(err)
	}

	// Set key.
	conn.Set(mykey, "Hello Redis")

	// Get key.
	val, err := conn.Get(mykey).ToString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)

	done := make(chan struct{}, 0)

	// Change Key in different connection.
	go func() {
		conn, err := client.Dial("")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		// Update key.
		if err = conn.Set(mykey, "Update mykey").Err(); err != nil {
			log.Fatal(err)
		}
		close(done)
	}()

	// Wait for go-routine.
	<-done
	// Wait for invalidation of key.
	<-invalidated

	// Output:
	// Hello Redis
	// Key invalidated
}
