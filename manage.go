// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

package cache

import (
    "sync"
    "time"

    "github.com/gomodule/redigo/redis"
)

var Manage Manager

type (
    // Manager
    // 管理器接口.
    Manager interface {
        AcquireConn() redis.Conn
        AcquireLocker(key string) Locker
        ReleaseLocker(x Locker)
    }

    // 管理器结构体.
    manager struct {
        poolConn   *redis.Pool
        poolLocker *sync.Pool
    }
)

// AcquireConn
// 获取Redis连接.
func (o *manager) AcquireConn() redis.Conn {
    return o.poolConn.Get()
}

// AcquireLocker
// 获取锁资源.
func (o *manager) AcquireLocker(key string) Locker {
    x := o.poolLocker.Get().(*locker)
    x.before(key)
    return x
}

// ReleaseLocker
// 使用结束释放回池.
func (o *manager) ReleaseLocker(x Locker) {
    x.(*locker).after()
    o.poolLocker.Put(x)
}

// 构造.
func (o *manager) init() *manager {
    o.poolLocker = &sync.Pool{New: func() interface{} { return (&locker{}).init() }}

    o.poolConn = &redis.Pool{
        MaxIdle:         Config.MaxIdle,
        MaxActive:       Config.MaxActive,
        IdleTimeout:     time.Duration(Config.IdleTimeout) * time.Second,
        Wait:            Config.Wait,
        MaxConnLifetime: time.Duration(Config.MaxConnLifetime) * time.Second,
    }

    o.poolConn.Dial = func() (redis.Conn, error) {
        return redis.Dial(Config.Network, Config.Address,
            redis.DialPassword(Config.Password),
            redis.DialDatabase(Config.Database),
            redis.DialConnectTimeout(time.Duration(Config.ConnectTimeout)*time.Second),
            redis.DialReadTimeout(time.Duration(Config.ReadTimeout)*time.Second),
            redis.DialWriteTimeout(time.Duration(Config.WriteTimeout)*time.Second),
            redis.DialKeepAlive(time.Duration(Config.MaxConnLifetime)*time.Second),
        )
    }

    return o
}
