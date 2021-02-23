// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package tests

import (
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/cache"
)

func TestLock(t *testing.T) {

	ctx := log.NewContext()

	lock := cache.NewLock("example")
	// lock.NotRenewal(ctx)

	str, err := lock.Set(ctx)
	if err != nil {
		t.Errorf("lock error: %v", err)
		return
	}

	defer lock.Unset(ctx, str)

	if str != "" {
		t.Logf("lock succeed: ")
		time.Sleep(time.Second * 30)
		return
	}

	t.Errorf("lock fail: ")


}
