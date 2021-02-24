// author: wsfuyibing <websearch@163.com>
// date: 2021-02-21

package cache

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/fuyibing/log/v2"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v3"
)

// Configuration struct.
type configuration struct {
	Address      string `yaml:"addr"`
	Database     int    `yaml:"database"`
	IdleTimeout  int    `yaml:"idle-timeout"`
	KeepAlive    int    `yaml:"keep-alive"`
	MaxActive    int    `yaml:"max-active"`
	MaxIdle      int    `yaml:"max-idle"`
	MaxLifetime  int    `yaml:"max-lifetime"`
	Network      string `yaml:"network"`
	Password     string `yaml:"password"`
	ReadTimeout  int    `yaml:"read-timeout"`
	Timeout      int    `yaml:"timeout"`
	Wait         bool   `yaml:"wait"`
	WriteTimeout int    `yaml:"write-timeout"`
	mutex        *sync.RWMutex
	pool         *redis.Pool
}

// Load configuration from specified yaml file.
func (o *configuration) LoadYaml(file string) error {
	// 1. read file content.
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// 2. assign parsed to instance.
	if err = yaml.Unmarshal(body, o); err != nil {
		return err
	}
	// 3.reset.
	log.Debugf("parse configuration from %s.", file)
	o.reset()
	return nil
}

// Return connection pool.
func (o *configuration) Pool() *redis.Pool { return o.pool }

// Initialize default configuration.
func (o *configuration) initialize() {
	o.mutex = new(sync.RWMutex)
	for _, file := range []string{"./tmp/cache.yaml", "./config/cache.yaml", "../config/cache.yaml"} {
		if nil == o.LoadYaml(file) {
			break
		}
	}
}

// Reset connection settings.
func (o *configuration) reset() {
	o.pool = &redis.Pool{MaxActive: o.MaxActive, MaxIdle: o.MaxIdle, Wait: o.Wait}
	// lifetime
	if o.MaxLifetime > 0 {
		o.pool.MaxConnLifetime = time.Duration(o.MaxLifetime) * time.Second
	}
	// timeout: idle
	if o.IdleTimeout > 0 {
		o.pool.IdleTimeout = time.Duration(o.IdleTimeout) * time.Second
	}
	// Connect
	o.pool.Dial = func() (redis.Conn, error) {
		// options: default.
		opts := make([]redis.DialOption, 0)
		opts = append(opts, redis.DialPassword(o.Password), redis.DialDatabase(o.Database))
		// options: timeouts.
		//          connect
		//          read
		//          write
		if o.Timeout > 0 {
			opts = append(opts, redis.DialConnectTimeout(time.Duration(o.Timeout)*time.Second))
		}
		if o.ReadTimeout > 0 {
			opts = append(opts, redis.DialReadTimeout(time.Duration(o.ReadTimeout)*time.Second))
		}
		if o.WriteTimeout > 0 {
			opts = append(opts, redis.DialWriteTimeout(time.Duration(o.WriteTimeout)*time.Second))
		}
		// options: keep alive
		if o.KeepAlive > 0 {
			opts = append(opts, redis.DialKeepAlive(time.Duration(o.KeepAlive)*time.Second))
		}
		// create connection
		return redis.Dial(o.Network, o.Address, opts...)
	}
}
