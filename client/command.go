// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

//go:generate commander

type command struct {
	send sendFct
}

func newCommand(send sendFct, sendInterceptor SendInterceptor) *command {
	if sendInterceptor == nil {
		return &command{send: send}
	}
	return &command{send: func(name string, r *result) {
		sendInterceptor(name, r.cmd())
		send(name, r)
	}}
}

// check interface implementations.
var _ Commands = (*command)(nil)
