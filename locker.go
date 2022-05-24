// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

package cache

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/fuyibing/log/v2"
    "github.com/gomodule/redigo/redis"
)

type (
    // Locker
    // 分布式锁接口.
    Locker interface {
        // Apply
        // 申请锁资源.
        Apply(ctx context.Context) (got bool, err error)

        // Release
        // 释放锁资源.
        Release()

        // Renewal
        // 续期状态.
        Renewal(renewal bool) Locker
    }

    // 分布式锁结构体.
    locker struct {
        mu                 *sync.RWMutex
        cancel             context.CancelFunc
        context, ctx       context.Context
        conn               redis.Conn
        key, value         string
        got, renewal       bool
        sending, listening bool
    }
)

// Apply
// 申请锁资源.
func (o *locker) Apply(ctx context.Context) (got bool, err error) {
    var reply interface{}

    // 1. 后置执行.
    defer func() {
        if r := recover(); r != nil {
            log.Panicfc(ctx, "[locker] apply fatal: %v.", err)
        } else {
            log.Infofc(ctx, "[locker] apply completed: key=%s, reply=%v.", o.key, reply)
        }
    }()

    // 2. 准备申请.
    o.ctx = ctx
    o.conn = Manage.AcquireConn()

    // 3. 申请资源.
    if log.Config.DebugOn() {
        log.Debugfc(ctx, "[locker] apply begin.")
    }
    if reply, err = o.do("SET", o.key, o.value, "NX", "EX", Config.LockerLifetime); err == nil {
        if got = fmt.Sprintf("%v", reply) == "OK"; got {
            o.got = got
            if o.renewal {
                o.listen()
            }
        }
    }
    return
}

// Release
// 释放锁资源.
func (o *locker) Release() {
    Manage.(*manager).releaseLocker(o)
}

// Renewal
// 续期状态.
func (o *locker) Renewal(renewal bool) Locker {
    o.renewal = renewal
    return o
}

// 后置.
func (o *locker) after() *locker {
    // 1. 退出上下文.
    if o.context.Err() == nil {
        o.cancel()
    }

    // 2. 等待结束.
    //    - listener
    //    - command
    for {
        if func() bool {
            o.mu.RLock()
            defer o.mu.RUnlock()
            return !o.sending && !o.listening
        }() {
            break
        }
    }

    // 3. 关闭连接.
    if o.conn != nil {
        // 删除资源.
        if err := o.conn.Send("DEL", o.key); err != nil {
            log.Warnfc(o.ctx, "[locker] delete applied resource error: %v.", err)
        } else if log.Config.DebugOn() {
            log.Debugfc(o.ctx, "[locker] delete applied resource.")
        }

        // 关闭连接.
        if err := o.conn.Close(); err != nil {
            log.Warnfc(o.ctx, "[locker] close connection error: %v.", err)
        } else if log.Config.DebugOn() {
            log.Debugfc(o.ctx, "[locker] close connection.")
        }

        // 重置字段.
        o.conn = nil
    }

    // 4. 重设上下文.
    o.cancel = nil
    o.context = nil

    // 5. 重置字段值.
    o.ctx = nil
    o.key = ""
    o.value = ""
    return o
}

// 前置.
func (o *locker) before(key string) *locker {
    // 1. 创建上下文.
    o.context, o.cancel = context.WithTimeout(context.Background(), time.Duration(Config.LockerRenewalTimeout)*time.Second)

    // 2. 字段初始化.
    o.key = fmt.Sprintf("%s:%s", Config.LockerPrefix, key)
    o.value = time.Now().String()
    o.got = false
    o.renewal = true
    o.sending = false
    o.listening = false
    return o
}

// 执行命令.
func (o *locker) do(cmd string, args ...interface{}) (reply interface{}, err error) {
    if o.context.Err() != nil {
        return
    }

    // 1. 阻塞命令.
    //    前一次命令未结束.
    if func() bool {
        o.mu.RLock()
        defer o.mu.RUnlock()
        return o.sending
    }() {
        return o.do(cmd, args...)
    }

    // 2. 更新状态.
    //    标记命令状态发送中.
    o.mu.Lock()
    o.sending = true
    o.mu.Unlock()

    // 3. 恢复状态.
    defer func() {
        o.mu.Lock()
        o.sending = false
        o.mu.Unlock()
    }()

    // 4. 发送命令.
    if log.Config.DebugOn() {
        log.Debugfc(o.ctx, "[locker] send command: %s, %v.", cmd, args)
    }
    return o.conn.Do(cmd, args...)
}

// 构造.
func (o *locker) init() *locker {
    o.mu = new(sync.RWMutex)
    return o
}

// 监听.
func (o *locker) listen() {
    // 1. 重复监听.
    if func() bool {
        o.mu.RLock()
        defer o.mu.RUnlock()
        return o.listening
    }() {
        return
    }

    // 2. 更新状态.
    o.mu.Lock()
    o.listening = true
    o.mu.Unlock()
    if log.Config.DebugOn() {
        log.Debugfc(o.ctx, "[locker] begin listen.")
    }

    // 3. 监控过程.
    go func() {
        // 结束监控.
        defer func() {
            // 捕获异常.
            if r := recover(); r != nil {
                log.Errorfc(o.ctx, "[locker] listen fatal: %v.", r)
            } else if log.Config.DebugOn() {
                log.Debugfc(o.ctx, "[locker] listen end.")
            }

            // 恢复状态.
            o.mu.Lock()
            o.listening = false
            o.mu.Unlock()
        }()

        // 定时续期.
        t := time.NewTicker(time.Duration(Config.LockerRenewalSeconds) * time.Second)
        for {
            select {
            case <-t.C:
                go o.update()
            case <-o.context.Done():
                return
            }
        }
    }()
}

// 续期.
func (o *locker) update() {
    // 1. 续期过程.
    reply, err := o.do("SET", o.key, o.value, "XX", "EX", Config.LockerLifetime)

    // 2. 续期出错.
    if err != nil {
        log.Errorfc(o.ctx, "[locker] renewal error: %v.", err)
        return
    }

    // 3. 续期完成.
    if log.Config.DebugOn() {
        log.Debugfc(o.ctx, "[locker] renewal: %v.", reply)
    }
}
