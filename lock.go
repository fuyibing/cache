// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var lockerPool *sync.Pool

type locker struct {
	cancel       context.CancelFunc
	context      context.Context
	conn         redis.Conn
	key, value   string
	got, renewal bool
	loser        chan bool
}

// NotRenewal
// 取消续期.
func (o *locker) NotRenewal(ctx interface{}) LockInterface {
	o.renewal = false
	return o
}

// Set
// 加锁.
func (o *locker) Set(ctx interface{}) (string, error) {
	var (
		err   error
		reply interface{}
	)

	// 1. 后置.
	defer func() {
		if err != nil {
			log.Errorfc(ctx, "[lock] apply locker resource error: key=%s, %v.", o.key, err)
		} else {
			log.Infofc(ctx, "[lock] apply locker resource: key=%s, reply=%v.", o.key, reply)
		}
	}()

	// 2. 申请过程.
	if reply, err = o.conn.Do("SET", o.key, o.value, "NX", "EX", LockLifetime); err != nil {
		return "", err
	}

	// 3. 校验结果.
	if fmt.Sprintf("%v", reply) == "OK" {
		o.got = true

		if o.renewal {
			o.listen(ctx)
		}

		return o.value, nil
	}

	return "", nil
}

// Unset
// 解锁.
func (o *locker) Unset(ctx interface{}, resource string) (err error) {
	var (
		reply interface{}
	)

	// 1. 删除资源.
	if o.got {
		if reply, err = o.conn.Do("DEL", o.key); err == nil {
			if log.Config.DebugOn() {
				log.Debugfc(ctx, "[locker] delete locked resource: %v.", reply)
			}
		} else {
			log.Errorfc(ctx, "[locker] delete locked resource error: %v.", err)
		}
	}

	// 2. 恢复字段.
	o.after(ctx)

	// 3. 释放回池.
	lockerPool.Put(o)
	return
}

// 后置.
func (o *locker) after(ctx interface{}) {
	// 1. 取消上下文.
	o.cancel()
	o.cancel = nil
	o.context = nil

	// 2. 关闭连接.
	if err := o.conn.Close(); err != nil {
		log.Errorfc(ctx, "[lock] close redis connection error: %v.", err)
	}
	o.conn = nil
	o.loser = nil

	// 3. 重置字段.
	o.key = ""
	o.value = ""
}

// 前置.
func (o *locker) before(key string) {
	o.context, o.cancel = context.WithCancel(context.Background())
	o.conn = Config.Pool().Get()
	o.got = false
	o.renewal = true
	o.loser = make(chan bool)
	o.key = fmt.Sprintf("%s:%s", LockPrefix, key)
	o.value = o.uuid()
}

// 构造.
func (o *locker) init() *locker {
	return o
}

// 监听.
// 加锁成功, 且开启续期时此方法被调用.
func (o *locker) listen(ctx interface{}) {
	// 1. 上下文退出.
	if o.context.Err() != nil {
		return
	}

	// 2. 监听过程.
	t := time.NewTicker(time.Duration(LockRenewal) * time.Second)
	if log.Config.DebugOn() {
		log.Debugfc(ctx, "[lock] start renewal listener.")
	}
	go func() {
		defer func() {
			if log.Config.DebugOn() {
				log.Debugfc(ctx, "[lock] stop renewal listener.")
			}
		}()
		for {
			select {
			case <-t.C:
				// 定时续期.
				go o.update(ctx)

			case <-o.loser:
				// 续期出错.
				return

			case <-o.context.Done():
				// 上下文退出.
				return
			}
		}
	}()
}

// 续期.
func (o *locker) update(ctx interface{}) {
	var (
		err   error
		reply interface{}
	)

	// 1. 后置.
	//    检查续期结果, 若续期失败则发送loser信号, 用于退出续期. 失败原因
	//    - 服务器错误
	//    - Key不存在
	defer func() {
		if err != nil {
			log.Infofc(ctx, "[lock] renewal error: %v.", err)
			o.loser <- true
		} else if log.Config.DebugOn() {
			log.Debugfc(ctx, "[lock] renewal completed.")
		}
	}()

	// 2. 续期过程.
	//    XX: Key必须存在
	//    EX: 续期时长
	if reply, err = o.conn.Do("SET", o.key, o.value, "XX", "EX", LockLifetime); err != nil {
		if fmt.Sprintf("%v", reply) != "OK" {
			err = fmt.Errorf("update failed")
		}
	}
}

// 生成键值.
func (o *locker) uuid() string {
	if u, e := uuid.NewUUID(); e == nil {
		return strings.ReplaceAll(u.String(), "-", "")
	}
	t := time.Now()
	return fmt.Sprintf("a%d%d%d", t.Unix(), t.UnixNano(), rand.Int63n(999999999999))
}

// NewLock
// return locker manager interface.
func NewLock(key string) LockInterface {
	o := lockerPool.Get().(*locker)
	o.before(key)
	return o
}
