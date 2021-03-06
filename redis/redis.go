package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	defaultPort        = 6379
	defaultMaxIdle     = 5
	defaultIdleTimeout = 300 * time.Second
	defaultNetwork     = "tcp"

	setIfNotExist     = "NX"
	setWithExpireTime = "PX"
	lockSuccess       = "OK"
	unlockSuccess     = 1
)

type RedisPool struct {
	pool        Pool
	redisConfig *RedisConfig
}

func NewRedisPool(host string, opts ...RedisOption) *RedisPool {
	config := &RedisConfig{
		network:     defaultNetwork,
		host:        host,
		port:        defaultPort,
		password:    "",
		db:          0,
		maxIdle:     defaultMaxIdle,
		idleTimeout: defaultIdleTimeout,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &RedisPool{
		pool: &redis.Pool{
			MaxIdle:     config.maxIdle,
			IdleTimeout: config.idleTimeout,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(config.network, fmt.Sprintf("%s:%d", config.host, config.port), redis.DialPassword(config.password), redis.DialDatabase(config.db))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		redisConfig: config,
	}
}

func (rp RedisPool) Ping() error {
	conn := rp.pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

/**
1. set value 加入了NX参数，可以保证如果已有key存在，则函数不会调用成功，也就是只有一个客户端能持有锁，满足互斥性;
2. 由于我们对锁设置了过期时间，即使锁的持有者后续发生崩溃而没有解锁，锁也会因为到了过期时间而自动解锁（即key被删除），不会发生死锁;
3. 将value=requestId，用来标识这把锁是属于哪个请求加的，那么在客户端在解锁的时候就可以进行校验是否是同一个客户端。
*/
func (rp RedisPool) Acquire(key, value string, expiry time.Duration) bool {
	conn := rp.pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("SET", key, value, setIfNotExist, setWithExpireTime, int(expiry/time.Millisecond)))
	return err == nil && reply == lockSuccess
}

const unlockScriptSrc = `
if redis.call('get', KEYS[1]) == ARGV[1] then 
	return redis.call('del', KEYS[1])
else
	return 0 end
`

/**
将Lua代码传到jedis.eval()方法里，并使参数KEYS[1]赋值为lockKey，ARGV[1]赋值为requestId。
在执行的时候，首先会获取锁对应的value值，检查是否与requestId相等，如果相等则解锁（删除key）。
*/
func (rp RedisPool) Release(key, value string) bool {
	unlockScript := redis.NewScript(1, unlockScriptSrc)
	conn := rp.pool.Get()
	defer conn.Close()
	reply, err := redis.Int64(unlockScript.Do(conn, key, value))
	return err == nil && reply == unlockSuccess
}
