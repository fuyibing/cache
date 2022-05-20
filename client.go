// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

package cache

import (
	"fmt"

	"github.com/fuyibing/log/v2"
	"github.com/gomodule/redigo/redis"
)

type (
	// Client
	// 客户端接口.
	Client interface {
		// Do
		// 执行命令.
		//
		// 此命令用于读取Redis数据, 执行过程中获取连接, 并在使用结束后关闭连接.
		Do(cmd string, args ...interface{}) (reply interface{}, err error)

		// DoWithConn
		// 执行命令.
		//
		// 此命令用于读取Redis数据.
		DoWithConn(conn redis.Conn, cmd string, args ...interface{}) (reply interface{}, err error)

		// Send
		// 执行命令.
		//
		// 此命令用于写入Redis数据, 执行过程中获取连接, 并在使用结束后关闭连接.
		Send(cmd string, args ...interface{}) (err error)

		// SendWithConn
		// 执行命令.
		//
		// 此命令用于写入Redis数据.
		SendWithConn(conn redis.Conn, cmd string, args ...interface{}) (err error)
	}

	// 客户端结构体.
	client struct{}
)

// Do
// 执行命令.
func (o *client) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := Manage.AcquireConn()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("[client] close connection error: %v.", err)
		}
	}()

	return o.DoWithConn(conn, cmd, args...)
}

// DoWithConn
// 执行命令.
func (o *client) DoWithConn(conn redis.Conn, cmd string, args ...interface{}) (reply interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return conn.Do(cmd, args...)
}

// Send
// 发送命令.
func (o *client) Send(cmd string, args ...interface{}) error {
	conn := Manage.AcquireConn()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("[client] close connection error: %v.", err)
		}
	}()

	return o.SendWithConn(conn, cmd, args...)
}

// SendWithConn
// 发送命令.
func (o *client) SendWithConn(conn redis.Conn, cmd string, args ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return conn.Send(cmd, args...)
}

// 构造.
func (o *client) init() *client {
	log.Info("[cache] initialize client.")
	return o
}
