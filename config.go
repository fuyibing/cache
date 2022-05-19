// author: wsfuyibing <websearch@163.com>
// date: 2022-05-19

package cache

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

const (
    defaultLockerLifetime       = 30
    defaultLockerPrefix         = "LOCK"
    defaultLockerRenewalSeconds = 5
    defaultLockerRenewalTimeout = 3600

    defaultMaxIdle         = 5
    defaultMaxActive       = 50
    defaultIdleTimeout     = 60
    defaultMaxConnLifetime = 60

    defaultConnectTimeout = 10
    defaultReadTimeout    = 5
    defaultWriteTimeout   = 5
)

// Config
// 配置实例.
var Config *Configuration

// Configuration
// 配置参数.
type Configuration struct {
    Network  string `yaml:"network"`
    Address  string `yaml:"address"`
    Database int    `yaml:"database"`
    Password string `yaml:"password"`

    ConnectTimeout int `yaml:"connect-timeout"`
    ReadTimeout    int `yaml:"read-timeout"`
    WriteTimeout   int `yaml:"write-timeout"`

    MaxIdle         int  `yaml:"max-idle"`
    MaxActive       int  `yaml:"max-active"`
    IdleTimeout     int  `yaml:"idle-timeout"`
    Wait            bool `yaml:"wait"`
    MaxConnLifetime int  `yaml:"max-conn-lifetime"`

    LockerLifetime       int    `yaml:"locker-lifetime"`
    LockerPrefix         string `yaml:"locker-prefix"`
    LockerRenewalSeconds int    `yaml:"locker-renewal-seconds"`
    LockerRenewalTimeout int    `yaml:"locker-renewal-timeout"`
}

// 赋值.
// 从YAML文件加载配置后, 未设置项赋默认值.
func (o *Configuration) defaults() {
    if o.MaxIdle < 1 {
        o.MaxIdle = defaultMaxIdle
    }
    if o.MaxActive < 1 {
        o.MaxActive = defaultMaxActive
    }
    if o.IdleTimeout < 1 {
        o.IdleTimeout = defaultIdleTimeout
    }
    if o.MaxConnLifetime < 1 {
        o.MaxConnLifetime = defaultMaxConnLifetime
    }

    if o.ConnectTimeout < 1 {
        o.ConnectTimeout = defaultConnectTimeout
    }
    if o.ReadTimeout < 1 {
        o.ReadTimeout = defaultReadTimeout
    }
    if o.WriteTimeout < 1 {
        o.WriteTimeout = defaultWriteTimeout
    }

    if o.LockerLifetime < 1 {
        o.LockerLifetime = defaultLockerLifetime
    }
    if o.LockerPrefix == "" {
        o.LockerPrefix = defaultLockerPrefix
    }
    if o.LockerRenewalSeconds < 1 {
        o.LockerRenewalSeconds = defaultLockerRenewalSeconds
    }
    if o.LockerRenewalTimeout < 1 {
        o.LockerRenewalTimeout = defaultLockerRenewalTimeout
    }

    fmt.Printf("config: %+v\n", o)
}

// 构造.
func (o *Configuration) init() *Configuration {
    o.load()
    o.defaults()
    return o
}

// 加载.
// 从YAML文件中加载配置参数.
func (o *Configuration) load() {
    for _, file := range []string{
        "../tmp/cache.yaml",
        "../config/cache.yaml",
        "./tmp/cache.yaml",
        "./config/cache.yaml",
    } {
        body, err := os.ReadFile(file)
        if err == nil {
            if err = yaml.Unmarshal(body, o); err == nil {
                break
            }
        }
    }
}
