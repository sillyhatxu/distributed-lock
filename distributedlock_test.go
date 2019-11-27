package distlock

import (
	"fmt"
	"github.com/sillyhatxu/distributed-lock/redis"
	"testing"
	"time"
)

const (
	testHost = "127.0.0.1"
	testPort = 16379
)

func TestNew(t *testing.T) {
	timeout := 60 * time.Second
	timeoutTime := time.Now().Add(timeout)
	fmt.Println(timeoutTime)
}

func TestLockGoroutine(t *testing.T) {
	t1 := time.Now()
	key := "testlockkey"
	threadNumber := 50
	number := 100
	expected := threadNumber * number
	count := 0
	redisPool := redis.NewRedisPool(testHost, redis.Port(testPort))
	redisLock, err := New(redisPool)
	if err != nil {
		panic(err)
	}
	//go testGoroutine(key, 1, number, redisLock)
	for i := 0; i < threadNumber; i++ {
		go func() {
			for i := 0; i < number; i++ {
				err := redisLock.Lock(key, func() error {
					count++
					return nil
				})
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "start")
	for count != expected {
		time.Sleep(10 * time.Second)
		fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "count:", count)
	}
	elapsed := time.Since(t1)
	fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "App elapsed: ", elapsed, " --- ", expected, count)

}
