// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/fuyibing/log/v2"
)

// Redis command client.
type client struct{}

// Delete key.
func (o *client) Del(ctx interface{}, keys ...interface{}) (Response, error) {
	return o.do(ctx, "DEL", keys...)
}

// Set expiration lifetime.
func (o *client) Expire(ctx interface{}, key string, seconds int) (Response, error) {
	return o.do(ctx, "EXPIRE", key, seconds)
}

// Read key.
func (o *client) Get(ctx interface{}, key string) (Response, error) {
	return o.do(ctx, "GET", key)
}

// Set key without lifetime.
func (o *client) Set(ctx interface{}, key string, value interface{}) (Response, error) {
	return o.do(ctx, "SET", key, value)
}

// Set key if not exist without lifetime.
func (o *client) SetNx(ctx interface{}, key string, value interface{}) (Response, error) {
	return o.do(ctx, "SET", key, value, "NX")
}

// Set key if not exist with lifetime.
func (o *client) SetNxEx(ctx interface{}, key string, value interface{}, seconds int) (Response, error) {
	return o.do(ctx, "SET", key, value, "NX", "EX", seconds)
}

// Run command.
func (o *client) Do(ctx interface{}, cmd string, args ...interface{}) (res Response, err error) {
	return o.do(ctx, cmd, args...)
}

// Run command.
func (o *client) do(ctx interface{}, cmd string, args ...interface{}) (res Response, err error) {
	t := time.Now()
	// 1. defer detect.
	defer func() {
		d := time.Now().Sub(t).Seconds()
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			log.Errorfc(ctx, "[cache][d=%f] run %s command got fatal error with arguments: %v: %v.", d, cmd, args, err)
		} else {
			if err != nil {
				log.Errorfc(ctx, "[cache][d=%f] run %s command error with arguments: %v: %v.", d, cmd, args, err)
			} else {
				log.Infofc(ctx, "[cache][d=%f] run %s command completed with arguments: %v.", d, cmd, args)
			}
		}
	}()
	// 2. Pool manager.
	if log.Config.DebugOn() {
		log.Debugfc(ctx, "[cache] acquire connection from pool.")
	}
	c := Config.Pool().Get()
	defer func() {
		if x := c.Close(); x != nil {
			log.Warnfc(ctx, "[cache] release connection error: %v.", x)
		} else {
			if log.Config.DebugOn() {
				log.Debugfc(ctx, "[cache] release connection to pool.")
			}
		}
	}()
	// 3. Send redis command.
	var v interface{}
	if v, err = c.Do(cmd, args...); err != nil {
		return
	}
	// 4. completed command.
	return &response{
		v: v,
	}, nil
}
