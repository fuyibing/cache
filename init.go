// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

// Package cache.
//
// 基于Redis的缓存管理, 支持Redis分步式锁.
package cache

import (
	"sync"

	"github.com/fuyibing/log/v2"
)

const (
	LockPrefix   = "LOCK"
	LockLifetime = 15
	LockRenewal  = 5
)

var (
	Config ConfigInterface
	Client ClientInterface
)

func init() {
	new(sync.Once).Do(func() {
		log.Debugf("init cache package.")
		Config = new(configuration)
		Config.initialize()
		Client = new(client)

		lockerPool = &sync.Pool{
			New: func() interface{} {
				return (&locker{}).init()
			},
		}
	})
}
