// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"log"
)

// GenericNotification represents the type for an out of bound push notification send by Redis.
type genericNotification struct {
	kind   string
	values []RedisValue
}

// SubscribeNotification represents the type for an out of bound subscribe push notification send by Redis.
type subscribeNotification struct {
	channel string // channel or pattern
	count   int64
}

// UnsubscribeNotification represents the type for an out of bound unsubscribe push notification send by Redis.
type unsubscribeNotification struct {
	channel string
	count   int64
}

// PublishNotification represents the type for an out of bound publish push notification send by Redis.
type publishNotification struct {
	pattern string
	channel string
	msg     string
}

// InvalidateNotification represents the type for an out of bound invalidation push notification send by Redis (client side caching).
type invalidateNotification struct {
	keys []string // keys to invalidate
}

//
const (
	pubSubSubscribe    = "subscribe"
	pubSubUnsubscribe  = "unsubscribe"
	pubSubPsubscribe   = "psubscribe"
	pubSubPunsubscribe = "punsubscribe"
	pubSubMessage      = "message"
	pubSubPMessage     = "pmessage"
	invalidateMessage  = "invalidate"
)

func assertNotification(condition bool, v []RedisValue) {
	if !condition {
		log.Panicf("invalid notification %v", v)
	}
}

//
func newNotification(v []RedisValue) (interface{}, error) {
	assertNotification(len(v) > 0 && v[0].Kind() == RkString, v)

	kind := string(v[0].(_string))

	switch kind {

	case pubSubSubscribe, pubSubPsubscribe:
		assertNotification(len(v) == 3 && v[1].Kind() == RkString && v[2].Kind() == RkNumber, v)
		return &subscribeNotification{channel: string(v[1].(_string)), count: int64(v[2].(_number))}, nil

	case pubSubUnsubscribe, pubSubPunsubscribe:
		assertNotification(len(v) == 3 && v[1].Kind() == RkString && v[2].Kind() == RkNumber, v)
		return &unsubscribeNotification{channel: string(v[1].(_string)), count: int64(v[2].(_number))}, nil

	case pubSubMessage:
		assertNotification(len(v) == 3 && v[1].Kind() == RkString && v[2].Kind() == RkString, v)
		return &publishNotification{channel: string(v[1].(_string)), msg: string(v[2].(_string))}, nil

	case pubSubPMessage:
		assertNotification(len(v) == 4 && v[1].Kind() == RkString && v[2].Kind() == RkString && v[3].Kind() == RkString, v)
		return &publishNotification{channel: string(v[1].(_string)), pattern: string(v[2].(_string)), msg: string(v[3].(_string))}, nil

	case invalidateMessage:
		assertNotification(len(v) == 2 && v[1].Kind() == RkSlice, v)
		keys, err := v[1].ToStringSlice()
		if err != nil {
			return nil, err
		}
		return &invalidateNotification{keys: keys}, nil

	default:
		return &genericNotification{kind: kind, values: v[1:]}, nil
	}
}
