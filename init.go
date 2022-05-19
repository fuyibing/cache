// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

// Package cache
// use to manage redis cache.
package cache

import "sync"

func init() {
    new(sync.Once).Do(func() {
        Config = (&Configuration{}).init()
        Manage = (&manager{}).init()
    })
}
