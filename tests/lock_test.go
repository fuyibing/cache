// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package tests

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cache"
)

func TestLock(t *testing.T) {
	println("goroutine: ", runtime.NumGoroutine())

	for i := 0; i < 1; i++ {
		testLock(t, i)
	}

	time.Sleep(time.Second)
	println("goroutine: ", runtime.NumGoroutine())
}

func testLock(t *testing.T, n int) {

	c := log.NewContext()
	x := cache.NewLock(fmt.Sprintf("example:%d", n))
	str, err := x.Set(c)
	if err != nil {
		t.Errorf("lock error: %v", err)
		return
	}

	if str == "" {
		return
	}

	defer func() {
		if err2 := x.Unset(c, str); err2 != nil {
			t.Logf("------------- unset error: %v", err2)
		}
	}()

	t.Logf("lock succeed: ")
	time.Sleep(time.Second)
}
