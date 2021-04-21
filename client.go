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

func (o *client) Decr(ctx interface{}, key string) (Response, error) {
	return o.Do(ctx, "DECR", key)
}

func (o *client) DecrBy(ctx interface{}, key string, step int) (Response, error) {
	return o.Do(ctx, "DECRBY", key, step)
}

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

func (o *client) Incr(ctx interface{}, key string) (Response, error) {
	return o.Do(ctx, "INCR", key)
}

func (o *client) IncrBy(ctx interface{}, key string, step int) (Response, error) {
	return o.Do(ctx, "INCRBY", key, step)
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

// ----------------
// HASH OPERATE
// ----------------

func (o *client) HGet(ctx interface{}, key, field string) (Response, error) {
	return o.Do(ctx, "HGET", key, field)
}

func (o *client) HGetAll(ctx interface{}, key string) (Response, error) {
	return o.Do(ctx, "HGETALL", key)
}

func (o *client) HDecr(ctx interface{}, key, field string) (Response, error) {
	return o.Do(ctx, "HDECR", key, field)
}

func (o *client) HDecrBy(ctx interface{}, key, field string, num interface{}) (Response, error) {
	return o.Do(ctx, "HDECRBY", key, field, num)
}

func (o *client) HIncr(ctx interface{}, key, field string) (Response, error) {
	return o.Do(ctx, "HINCR", key, field)
}

func (o *client) HIncrBy(ctx interface{}, key, field string, num interface{}) (Response, error) {
	return o.Do(ctx, "HINCRBY", key, field, num)
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
			log.Panicfc(ctx, "[cache] fatal error: %v.", r)
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
			log.Panicfc(ctx, "[cache][d=%f] fatal error: %v.", r)
			err = errors.New(fmt.Sprintf("%v", r))
		}
		if err != nil {
			log.Errorfc(ctx, "[cache][d=%f] send command error: %s %v: %s.", d1, cmd, args, err)
		} else {
			log.Infofc(ctx, "[cache][d=%f] send command completed: %s %v.", d1, cmd, args)
		}
	}()
	// 3.1 Send commend.
	var value interface{}
	if value, err = connection.Do(cmd, args...); err != nil {
		return
	}
	// 4. response.
	key := ""
	if len(args) > 0 {
		key = fmt.Sprintf("%v", args[0])
	}
	res = &response{exist: value != nil, key: key, value: value}
	return
}
