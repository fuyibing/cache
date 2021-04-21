// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"reflect"

	"github.com/gomodule/redigo/redis"
)

// Config manager.
type ConfigInterface interface {
	Pool() *redis.Pool
	initialize()
}

// Lock manager.
type LockInterface interface {
	// Not use lifetime renewal.
	NotRenewal(ctx interface{}) LockInterface

	// Lock resource.
	Set(ctx interface{}) (value string, err error)

	// Release locked resource.
	Unset(ctx interface{}, resource string) error
}

// Redis command client.
type ClientInterface interface {
	// Decrement 1.
	Decr(ctx interface{}, key string) (Response, error)

	// Decrement with step.
	DecrBy(ctx interface{}, key string, step int) (Response, error)

	// Delete key.
	Del(ctx interface{}, keys ...interface{}) (Response, error)

	// Do any command with specified arguments.
	Do(ctx interface{}, cmd string, args ...interface{}) (res Response, err error)

	// Set expiration lifetime.
	Expire(ctx interface{}, key string, seconds int) (Response, error)

	// Read key.
	Get(ctx interface{}, key string) (Response, error)

	// Increment with 1.
	Incr(ctx interface{}, key string) (Response, error)

	// Increment with step.
	IncrBy(ctx interface{}, key string, step int) (Response, error)

	// Set key without lifetime.
	Set(ctx interface{}, key string, value interface{}) (Response, error)

	// Set key if not exist without lifetime.
	SetNx(ctx interface{}, key string, value interface{}) (Response, error)

	// Set key if not exist with lifetime.
	SetNxEx(ctx interface{}, key string, value interface{}, seconds int) (Response, error)
}

// Response of redis command done.
// Contains read and write command.
type Response interface {

	//
	Exist() bool

	// Return value is equal to specified int.
	EqInt(i int) bool

	// Return value is equal to specified int32.
	EqInt32(i int32) bool

	// Return value is equal to specified int64.
	EqInt64(i int64) bool

	// Return value is equal to specified string.
	EqString(string) bool

	// Return int value.
	Int() int

	// Return int64 value.
	Int64() int64

	// Return value is nil.
	IsNil() bool

	// Return value is "OK" string.
	IsOk() bool

	// Return redis key name.
	Key() string

	// Return string for any response.
	String() string

	// Return reflection type.
	Type() reflect.Type

	// Return origin value.
	Value() interface{}
}
