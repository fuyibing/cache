// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cache"
)

func TestPool(t *testing.T) {
	wg := new(sync.WaitGroup)
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go run(wg, i)
	}
	wg.Wait()
	log.Infof("IdleCount: %d.", cache.Config.Pool().Stats().IdleCount)
	log.Infof("ActiveCount: %d.", cache.Config.Pool().Stats().ActiveCount)
	log.Infof("WaitCount: %d.", cache.Config.Pool().Stats().WaitCount)
	log.Infof("WaitDuration: %d.", cache.Config.Pool().Stats().WaitDuration)
	// time.Sleep(time.Second * 5)
}

func TestPoolLock(t *testing.T) {
	wg := new(sync.WaitGroup)
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go lock(wg, i)
	}
	wg.Wait()
	log.Infof("IdleCount: %d.", cache.Config.Pool().Stats().IdleCount)
	log.Infof("ActiveCount: %d.", cache.Config.Pool().Stats().ActiveCount)
	log.Infof("WaitCount: %d.", cache.Config.Pool().Stats().WaitCount)
	log.Infof("WaitDuration: %d.", cache.Config.Pool().Stats().WaitDuration)
	// time.Sleep(time.Second * 5)
}

func run(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	ctx := log.NewContext()
	cache.Client.Set(ctx, fmt.Sprintf("key:%d", i), i)
}

func lock(wg *sync.WaitGroup, i int) {
	defer func() {
		wg.Done()
	}()
	l := cache.NewLock(fmt.Sprintf("key:%d", i))
	k, err := l.Set(nil)
	if err == nil {
		time.Sleep(time.Millisecond)
		_ = l.Unset(nil, k)
	} else {
		log.Warnf("lock fail: %v.", err)
	}
}
