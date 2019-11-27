package distlock

import (
	"fmt"
	"github.com/sillyhatxu/distributed-lock/redis"
	"github.com/sirupsen/logrus"
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
	err := pool.Ping()
	return &DistributedLock{
		pool:          pool,
		config:        config,
		customChannel: make(chan error),
	}, err
}

type ExecuteFunc func() error

func (dl *DistributedLock) Lock(lockKey string, executeFun ExecuteFunc) error {
	if dl == nil || dl.config == nil {
		return fmt.Errorf("redis lock is nil")
	}
	key := dl.config.lockKey(lockKey)
	requestId := GeneratorRequestId()
	go dl.execute(key, requestId, executeFun, dl.customChannel)
	err := <-dl.customChannel
	dl.pool.Release(key, requestId)
	return err
}

func (dl *DistributedLock) execute(key string, requestId string, executeFun ExecuteFunc, c chan error) {
	var n uint
	for n < dl.config.attempts {
		if dl.pool.Acquire(key, requestId, dl.config.expiry) {
			c <- executeFun()
			return
		}
		if n >= dl.config.attempts-1 {
			break
		}
		time.Sleep(dl.config.delayType(n, dl.config))
		n++
		continue
	}
	logrus.Errorf("more than the number of retries : %v", ErrFailed)
	c <- ErrFailed
	return
}
