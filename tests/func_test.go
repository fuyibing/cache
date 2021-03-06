// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/fuyibing/cache"
)

func TestFunc(t *testing.T) {
	res, err := cache.Client.Get(nil, "test")
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		if res.EqString("name") {
			t.Logf("type: %v", res.Type().String())
			t.Logf("succeed: %s", res.Value())
		} else {
			t.Logf("failed: %s", res.Value())
		}
	}
	time.Sleep(time.Second)
}

func TestFuncDel(t *testing.T) {
	res, err := cache.Client.Del(nil, "test", "test1", "test2")
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		// if res.EqString("name") {
		t.Logf("type: %v", res.Type().String())
		t.Logf("succeed: %v", res.Value())
		// } else {
		// 	t.Logf("failed: %s", res.Value())
		// }
	}
	time.Sleep(time.Second)
}

func TestFuncExpire(t *testing.T) {
	res, err := cache.Client.Expire(nil, "test", 10)
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		if res.EqInt64(1) {
			t.Logf("result: Succeed")
		} else {
			t.Error("result: Failed")
		}
	}
	time.Sleep(time.Second)
}

func TestFuncGet(t *testing.T) {
	res, err := cache.Client.Get(nil, "key:name2")
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		t.Logf("key: %s.", res.Key())
		t.Logf("found: %v.", res.Exist())
		t.Logf("succeed: %v", res.String())
	}
	time.Sleep(time.Second)
}

func TestFuncSetNx(t *testing.T) {

	res, err := cache.Client.SetNx(nil, "test2", "value")
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		t.Logf("result: %v", res.Value())
	}

	time.Sleep(time.Second)
}

func TestFuncSetNxEx(t *testing.T) {

	res, err := cache.Client.SetNxEx(nil, "test3", "value", 5)
	if err != nil {
		t.Errorf("error: %v.", err)
	} else {
		if res.IsOk() {
			t.Logf("result: Succeed")
		} else {
			t.Logf("result: failed")
		}
	}

	time.Sleep(time.Second)
}

func TestFuncInre(t *testing.T) {

	r1, _ := cache.Client.Incr(nil, "test-1")
	t.Logf("r1.value: %v == 1.", r1.Int64())

	r2, _ := cache.Client.IncrBy(nil, "test-1", 3)
	t.Logf("r2.value: %v == 4.", r2.Int64())

	r3, _ := cache.Client.Decr(nil, "test-1")
	t.Logf("r3.value: %v == 3.", r3.Int64())

	r4, _ := cache.Client.DecrBy(nil, "test-1", 3)
	t.Logf("r4.value: %v == 0.", r4.Int64())
}


func TestFuncIncrInt(t *testing.T) {

	r1, _ := cache.Client.Incr(nil, "test-1")
	t.Logf("r1.value: %v.", r1.Int())
	t.Logf("r1.value: %v.", r1.Int64())

	x := r1.Int()

	t.Logf("x reflect: %s.", reflect.TypeOf(x).Kind().String())




}
