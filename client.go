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
	// 1. Panic.
	defer func() {
		if r := recover(); r != nil {
			log.Errorfc(ctx, "[cache] fatal error: %v.", r)
		}
	}()
	// 2. Connection.
	//    Acquire from pool before command and
	//    release when completed.
	t0 := time.Now()
	connection := Config.Pool().Get()
	if log.Config.DebugOn() {
		log.Debugfc(ctx, "[cache][d=%f] connection acquired from pool.", time.Now().Sub(t0).Seconds())
	}
	defer func() {
		if e0 := connection.Close(); e0 != nil {
			log.Errorfc(ctx, "[cache] release connection error: %v.", e0)
		} else {
			if log.Config.DebugOn() {
				log.Debugfc(ctx, "[cache] release connection to pool.")
			}
		}
	}()
	// 3. Command.
	//    Send redis command with arguments to server and
	//    generate response struct.
	t1 := time.Now()
	defer func() {
		d1 := time.Now().Sub(t1).Seconds()
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		if err != nil {
			log.Errorfc(ctx, "[cache][d=%f] send command error: %s %v: %s.", d1, cmd, args, err)
		} else {
			log.Infofc(ctx, "[cache][d=%f] send command completed: %s %v.", d1, cmd, args)
		}
	}()
	// 3.1 Send commend.
	var v interface{}
	if v, err = connection.Do(cmd, args...); err == nil {
		res = &response{v: v}
	}
	return
}
