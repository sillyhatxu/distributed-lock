package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// A Pool maintains a pool of Redis connections.
type Pool interface {
	Get() redis.Conn
}

type RedisConfig struct {
	network     string
	host        string
	port        int
	password    string
	db          int
	maxIdle     int
	idleTimeout time.Duration
}

type RedisOption func(*RedisConfig)

func Network(network string) RedisOption {
	return func(c *RedisConfig) {
		c.network = network
	}
}

func Host(host string) RedisOption {
	return func(c *RedisConfig) {
		c.host = host
	}
}

func Port(port int) RedisOption {
	return func(c *RedisConfig) {
		c.port = port
	}
}

func Password(password string) RedisOption {
	return func(c *RedisConfig) {
		c.password = password
	}
}

func DB(db int) RedisOption {
	return func(c *RedisConfig) {
		c.db = db
	}
}

func MaxIdle(maxIdle int) RedisOption {
	return func(c *RedisConfig) {
		c.maxIdle = maxIdle
	}
}

func IdleTimeout(idleTimeout time.Duration) RedisOption {
	return func(c *RedisConfig) {
		c.idleTimeout = idleTimeout
	}
}
