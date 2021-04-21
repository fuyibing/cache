// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"fmt"
	"reflect"
	"strconv"
)

// Response struct.
type response struct {
	exist bool
	key   string
	value interface{}
	ref   reflect.Type
}

// Exist in redis or not.
func (o *response) Exist() bool {
	return o.exist
}

// Return value is equal to specified int.
func (o *response) EqInt(i int) bool { return o.EqInt64(int64(i)) }

// Return value is equal to specified int32.
func (o *response) EqInt32(i int32) bool { return o.EqInt64(int64(i)) }

// Return value is equal to specified int64.
func (o *response) EqInt64(i int64) bool {
	if o.value != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				return o.value.(int64) == i
			case reflect.String:
				if n, err := strconv.ParseInt(o.value.(string), 0, 64); err == nil {
					return n == i
				}
			}
		}
	}
	return false
}

// Return value is equal to specified string.
func (o *response) EqString(str string) bool {
	if o.value != nil {
		return str == fmt.Sprintf("%s", o.value)
	}
	return false
}

func (o *response) Int() int {
	if o.value != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				if n, err := strconv.ParseInt(fmt.Sprintf("%v", o.value), 0, 32); err == nil {
					return int(n)
				}
			case reflect.String:
				if n, err := strconv.ParseInt(o.value.(string), 0, 64); err == nil {
					return int(n)
				}
			}
		}
	}
	return 0
}

func (o *response) Int64() int64 {
	if o.value != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				return o.value.(int64)
			case reflect.String:
				if n, err := strconv.ParseInt(o.value.(string), 0, 64); err == nil {
					return int64(n)
				}
			}
		}
	}
	return 0
}

// Return value is nil.
func (o *response) IsNil() bool {
	return o.value == nil
}

// Return value is "OK" string.
func (o *response) IsOk() bool {
	return o.EqString("OK")
}

// Return redis key name.
func (o *response) Key() string {
	return o.key
}

// Return string for value.
func (o *response) String() (str string) {
	// return empty string if value is nil.
	if o.value == nil {
		return
	}
	// return string dependent on reflection.
	t := o.t()
	switch t.Kind() {
	case reflect.String:
		str = o.value.(string)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = fmt.Sprintf("%d", o.value)
	case reflect.Slice:
		str = fmt.Sprintf("%s", o.value)
	default:
		str = fmt.Sprintf("%v", o.value)
	}
	return
}

// Return reflection type.
// Returned type equal to t() method.
func (o *response) Type() reflect.Type {
	return o.t()
}

// Return origin value.
// Returned value is dependent on redis stored.
func (o *response) Value() interface{} {
	return o.value
}

// Reflection type.
// Returned is nil or reflection type.
func (o *response) t() reflect.Type {
	// 1. return nil if not exist.
	if o.value == nil {
		return nil
	}
	// 2. reflect if not TypeOf called.
	if o.ref == nil {
		o.ref = reflect.TypeOf(o.value)
	}
	return o.ref
}
