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
	defaultTimeout       = 60 * time.Second
	defaultAttempts      = 8
	defaultDelay         = 10 * time.Millisecond
)

type ExecutionInterface interface {
	Exe() error
}

type DistributedLock struct {
	pool          *redis.RedisPool
	config        *ConfigOption
	customChannel chan ChannelResult
}

type ChannelResult struct {
	key       string
	requestId string
	err       error
}

func New(pool *redis.RedisPool, opts ...Option) (*DistributedLock, error) {
	//default
	config := &ConfigOption{
		lockKeyPrefix: defaultLockKeyPrefix,
		expiry:        defaultExpiry,
		timeout:       defaultTimeout,
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
		customChannel: make(chan ChannelResult),
	}, err
}

type ExecuteFunc func() error

func (dl *DistributedLock) Lock(lockKey string, execute ExecutionInterface) error {
	if dl == nil || dl.config == nil {
		return fmt.Errorf("redis lock is nil")
	}
	key := dl.config.lockKey(lockKey)
	requestId := GeneratorRequestId()
	go dl.execute(key, requestId, execute, dl.customChannel)
	channelResult := <-dl.customChannel
	if channelResult.err != nil {
		return channelResult.err
	}
	if !dl.pool.Release(channelResult.key, channelResult.requestId) {
		return fmt.Errorf("release error. %s : %s", key, requestId)
	}
	return nil
}

func (dl *DistributedLock) execute(key string, requestId string, execute ExecutionInterface, c chan ChannelResult) {
	timeoutTime := time.Now().Add(dl.config.timeout)
	for {
		if dl.pool.Acquire(key, requestId, dl.config.expiry) {
			c <- ChannelResult{key: key, requestId: requestId, err: execute.Exe()}
			return
		}
		if time.Now().UnixNano() >= timeoutTime.UnixNano() {
			fmt.Println(time.Now().UnixNano(), timeoutTime.UnixNano())
			break
		}
		time.Sleep(dl.config.delayType(1, dl.config))
		continue
	}
	logrus.Errorf("more than the number of retries : %v", ErrFailed)
	c <- ChannelResult{key: key, requestId: requestId, err: ErrFailed}
	return
}
