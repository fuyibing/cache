// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"fmt"
	"reflect"
	"strconv"
)

type response struct {
	v interface{}
}

// Return value is equal to specified int.
func (o *response) EqInt(i int) bool { return o.EqInt64(int64(i)) }

// Return value is equal to specified int32.
func (o *response) EqInt32(i int32) bool { return o.EqInt64(int64(i)) }

// Return value is equal to specified int64.
func (o *response) EqInt64(i int64) bool {
	if o.v != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				return o.v.(int64) == i
			case reflect.String:
				if n, err := strconv.ParseInt(o.v.(string), 0, 64); err == nil {
					return n == i
				}
			}
		}
	}
	return false
}

// Return value is equal to specified string.
func (o *response) EqString(str string) bool {
	if o.v != nil {
		return str == fmt.Sprintf("%s", o.v)
	}
	return false
}

func (o *response) Int() int {
	if o.v != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				if n, err := strconv.ParseInt(fmt.Sprintf("%v", o.v), 0, 32); err == nil {
					return int(n)
				}
			case reflect.String:
				if n, err := strconv.ParseInt(o.v.(string), 0, 64); err == nil {
					return int(n)
				}
			}
		}
	}
	return 0
}

func (o *response) Int64() int64 {
	if o.v != nil {
		if t := o.t(); t != nil {
			switch t.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int32, reflect.Int64:
				return o.v.(int64)
			case reflect.String:
				if n, err := strconv.ParseInt(o.v.(string), 0, 64); err == nil {
					return int64(n)
				}
			}
		}
	}
	return 0
}

// Return value is nil.
func (o *response) IsNil() bool {
	return o.v == nil
}

// Return value is "OK" string.
func (o *response) IsOk() bool {
	return o.EqString("OK")
}

// Return string for value.
func (o *response) String() string {
	return fmt.Sprintf("%v", o.v)
}

// Return reflection type.
func (o *response) Type() reflect.Type { return o.t() }

// Return origin value.
func (o *response) Value() interface{} { return o.v }

// Return reflection type.
func (o *response) t() reflect.Type { return reflect.TypeOf(o.v) }
