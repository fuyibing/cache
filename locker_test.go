// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

package cache

import (
	"testing"
	"time"

	"github.com/fuyibing/log/v2"
)

func TestManager_AcquireLocker(t *testing.T) {
	cli := Manage.AcquireLocker("example:100")
	defer func() {
		cli.Release()
	}()

	ctx := log.NewContext()
	log.Infofc(ctx, "locker ready.")

	got, err := cli.Apply(ctx)
	t.Logf("apply: %v->%v.", got, err)

	if got {
		time.Sleep(time.Second * 15)
	}
	t.Logf("end locker")
}
