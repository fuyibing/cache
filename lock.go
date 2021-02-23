// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
	"github.com/google/uuid"
)

// Lock struct.
type lock struct {
	ch        chan bool
	key       string
	listening bool
	mutex     *sync.RWMutex
	receipt   string
	renewal   bool
}

// Not use lifetime renewal.
func (o *lock) NotRenewal(ctx interface{}) LockInterface {
	log.Debugfc(ctx, "[cache] not use lifetime renewal.")
	o.renewal = false
	return o
}

// Lock resource.
func (o *lock) Set(ctx interface{}) (string, error) {
	v := o.uuid()
	// set redis.
	res, err := Client.SetNxEx(ctx, o.key, v, LockLifetime)
	if err != nil {
		return "", err
	}
	// res check.
	if res.IsOk() {
		if o.renewal {
			o.listen(ctx)
		}
		return v, nil
	}
	return "", nil
}

// Release locked resource.
func (o *lock) Unset(ctx interface{}, value string) error {
	// send quit channel if listening.
	if o.listening {
		o.ch <- true
	}
	// get locked resource.
	res, err := Client.Get(ctx, o.key)
	if err != nil {
		return err
	}
	// delete already if not exist.
	if res.IsNil() {
		return nil
	}
	// can not delete if value not equal to receipt value.
	if !res.EqString(value) {
		return errors.New(fmt.Sprintf("access denied"))
	}
	// send delete command.
	if _, err = Client.Del(ctx, o.key); err != nil {
		return err
	}
	return nil
}

// Set lifetime renewal.
func (o *lock) expiration(ctx interface{}) {
	_, _ = Client.Expire(ctx, o.key, LockLifetime)
}

// Listen channel.
func (o *lock) listen(ctx interface{}) {
	go func() {
		o.listening = true
		log.Debugfc(ctx, "[cache] listen lifetime renewal.")
		t := time.NewTicker(time.Duration(LockRenewal) * time.Second)
		defer func() {
			t.Stop()
			o.listening = false
		}()
		// listen channel
		for {
			select {
			case <-t.C:
				go o.expiration(ctx)
			case <-o.ch:
				return
			}
		}
	}()
}

// Generate unique identify string.
func (o *lock) uuid() string {
	if u, e := uuid.NewUUID(); e == nil {
		return strings.ReplaceAll(u.String(), "-", "")
	}
	t := time.Now()
	return fmt.Sprintf("a%d%d%d", t.Unix(), t.UnixNano(), rand.Int63n(999999999999))
}

// New lock instance.
func NewLock(key string) LockInterface {
	o := &lock{key: fmt.Sprintf("%s:%s", LockPrefix, key), listening: false, mutex: new(sync.RWMutex), renewal: true}
	o.ch = make(chan bool)
	return o
}
