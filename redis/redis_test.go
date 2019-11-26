package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	host     = "127.0.0.1"
	port     = 16379
	database = 4
	password = ""
)

func TestNewPool(t *testing.T) {
	redisPool := NewRedisPool(host, Port(port))
	err := redisPool.Ping()
	assert.Nil(t, err)
}

func TestAcquire(t *testing.T) {
	key := "lockkey"
	value := "this is value"
	expiry := 60 * time.Second
	redisPool := NewRedisPool(host, Port(port))
	result := redisPool.Acquire(key, value, expiry)
	assert.EqualValues(t, true, result)
}

func TestScriptGetValue(t *testing.T) {
	name := "dGVzdC1rZXk="
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()
	fmt.Println("conn : ", conn)
	test, err := redis.String(conn.Do("GET", name))
	fmt.Println("err is nil : ", err)
	fmt.Println("get name : ", test)
}

func TestScriptDeleteValue(t *testing.T) {
	name := "test-name"
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()
	fmt.Println("conn : ", conn)
	test, err := redis.String(conn.Do("DELETE", name))
	fmt.Println("err is : ", err)
	fmt.Println("get name : ", test)
}

func TestScriptExpire(t *testing.T) {
	name := "test-name"
	expiry := 5 * time.Second
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()
	expireScript := redis.NewScript(1, `
		redis.call("pexpire", KEYS[1], ARGV[1])
	`)
	result, err := redis.Int64(expireScript.Do(conn, name, int(expiry/time.Millisecond)))
	fmt.Println("err is : ", err)
	fmt.Println("result : ", result)
}

func TestScriptIncr(t *testing.T) {
	name := "test-name"
	quorum := 1
	expiry := 2 * time.Second
	count, expected := 100, 0
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()

	for i := 0; i < count; i++ {
		incrScript := redis.NewScript(1, `
		local current = redis.call("incr",KEYS[1])
		redis.call("pexpire", KEYS[1], ARGV[2])
		return current
		`)
		result, err := redis.Int64(incrScript.Do(conn, name, quorum, int(expiry/time.Millisecond)))
		if err != nil {
			panic(err)
		}
		fmt.Println("result : ", result)
		time.Sleep(100 * time.Millisecond)
		expected++
	}
	assert.EqualValues(t, expected, count)
}

func TestScriptDecr(t *testing.T) {
	name := "test-name"
	quorum := 1
	expiry := 2 * time.Second
	count, expected := 100, 0
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()

	for i := 0; i < count; i++ {
		incrScript := redis.NewScript(1, `
		return redis.call("decr",KEYS[1])
		`)
		result, err := redis.Int64(incrScript.Do(conn, name, quorum, int(expiry/time.Millisecond)))
		if err != nil {
			panic(err)
		}
		fmt.Println("result : ", result)
		time.Sleep(100 * time.Millisecond)
		expected++
	}
	assert.EqualValues(t, expected, count)
}

func TestScriptLock(t *testing.T) {
	name := "test-name"
	expiry := 2 * time.Second
	count, expected := 100, 0
	pool := getPool()
	conn := pool.Get()
	assert.Nil(t, conn.Err())
	defer conn.Close()

	var incrScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return 'LOCK'
	else
		local current = redis.call("SET",KEYS[1], ARGV[1])
		redis.call("pexpire", KEYS[1], ARGV[2])
		return current
	end
	`)
	//var incrScript = redis.NewScript(1, `
	//if redis.call("GET", KEYS[1]) == ARGV[1] then
	//	local current = redis.call("incr",KEYS[1])
	//	redis.call("pexpire", KEYS[1], ARGV[1])
	//	return current
	//	return redis.call("pexpire", KEYS[1], ARGV[2])
	//else
	//	return 0
	//end
	//`)

	for i := 0; i < count; i++ {
		result, err := redis.String(incrScript.Do(conn, name, "LOCK", int(expiry/time.Millisecond)))
		if err != nil {
			panic(err)
		}
		fmt.Println("result : ", result)
		time.Sleep(100 * time.Millisecond)
		expected++
	}
	assert.EqualValues(t, expected, count)
}

func getPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1", redis.DialPassword(password), redis.DialDatabase(database))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
