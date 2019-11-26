package distlock

import (
	"encoding/base64"
	"time"
)

type DelayTypeFunc func(n uint, co *ConfigOption) time.Duration

type ErrorCallbackFunc func(n uint, err error)

type LockKeyFunc func(key string) string

type ConfigOption struct {
	lockKeyPrefix string

	attempts      uint
	delay         time.Duration
	errorCallback ErrorCallbackFunc
	delayType     DelayTypeFunc

	lockKey LockKeyFunc
}

type Option func(*ConfigOption)

func LockKeyPrefix(lockKeyPrefix string) Option {
	return func(c *ConfigOption) {
		c.lockKeyPrefix = lockKeyPrefix
	}
}

//func Expiry(expiry time.Duration) Option {
//	return func(c *ConfigOption) {
//		c.expiry = expiry
//	}
//}
//
//func Attempts(attempts uint) Option {
//	return func(c *ConfigOption) {
//		c.attempts = attempts
//	}
//}
//
//func Delay(delay time.Duration) Option {
//	return func(c *ConfigOption) {
//		c.delay = delay
//	}
//}
//
//func DelayType(delayType DelayTypeFunc) Option {
//	return func(c *ConfigOption) {
//		c.delayType = delayType
//	}
//}
//
//func ErrorCallback(errorCallbackFunc ErrorCallbackFunc) Option {
//	return func(c *ConfigOption) {
//		c.errorCallback = errorCallbackFunc
//	}
//}
//
//func Network(network string) Option {
//	return func(c *ConfigOption) {
//		c.network = network
//	}
//}
//
//func Host(host string) Option {
//	return func(c *ConfigOption) {
//		c.host = host
//	}
//}
//
//func Port(port int) Option {
//	return func(c *ConfigOption) {
//		c.port = port
//	}
//}
//
//func Password(password string) Option {
//	return func(c *ConfigOption) {
//		c.password = password
//	}
//}
//
//func DB(db int) Option {
//	return func(c *ConfigOption) {
//		c.db = db
//	}
//}
//
//func MaxIdle(maxIdle int) Option {
//	return func(c *ConfigOption) {
//		c.maxIdle = maxIdle
//	}
//}
//
//func IdleTimeout(idleTimeout time.Duration) Option {
//	return func(c *ConfigOption) {
//		c.idleTimeout = idleTimeout
//	}
//}

func LockKey(lockKeyFunc LockKeyFunc) Option {
	return func(c *ConfigOption) {
		c.lockKey = lockKeyFunc
	}
}

func BackOffDelay(n uint, co *ConfigOption) time.Duration {
	result := co.delay * (1 << n)
	if result > 2*time.Second {
		return time.Second
	}
	return result
	//return co.delay * (1 << n)
}

func GeneratorLockKey(lockKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(lockKey))
}
