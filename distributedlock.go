package distlock

import (
	"github.com/sillyhatxu/distributed-lock/redis"
	"time"
)

const (
	defaultLockKeyPrefix = "GoDistRL" //golang distributed redis lock
	defaultExpiry        = 60 * time.Second
	defaultAttempts      = 8
	defaultDelay         = 200 * time.Millisecond
)

type DistributedLock struct {
	pool          *redis.RedisPool
	config        *ConfigOption
	customChannel chan error
}

func New(pool *redis.RedisPool, opts ...Option) (*DistributedLock, error) {
	//default
	config := &ConfigOption{
		lockKeyPrefix: defaultLockKeyPrefix,
		expiry:        defaultExpiry,
		attempts:      defaultAttempts,
		delay:         defaultDelay,
		errorCallback: func(n uint, err error) {

		},
		lockKey:   GeneratorLockKey,
		delayType: BackOffDelay,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &DistributedLock{
		pool:          pool,
		config:        config,
		customChannel: make(chan error),
	}, nil
}
