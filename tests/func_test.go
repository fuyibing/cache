// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package tests

import (
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
		t.Logf("type: %v", res.Type())
		if res.EqInt64(1) {
			t.Logf("result: Succeed")
		} else {
			t.Error("result: Failed")
		}
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
