package distlock

import (
	"fmt"
	"github.com/sillyhatxu/distributed-lock/redis"
	"sync"
	"testing"
	"time"
)

const (
	testHost = "127.0.0.1"
	testPort = 6379
)

type ParametersInterface struct {
	Params int
}

func (pi *ParametersInterface) Exe() error {
	//fmt.Printf("index: %d \n", pi.Params)
	fmt.Sprintf("index: %d", pi.Params)
	//pi.Params++
	return nil
}

func TestNew(t *testing.T) {
	timeout := 60 * time.Second
	timeoutTime := time.Now().Add(timeout)
	fmt.Println(timeoutTime)
}

//13.268109448s
//21.999689916s
//26.270671563s
//10.616449813s
func TestLockGoroutine(t *testing.T) {
	var wg sync.WaitGroup
	t1 := time.Now()
	key := "testlockkey1"
	threadNumber := 500
	number := 10
	expected := threadNumber * number
	count := 0
	redisPool := redis.NewRedisPool(testHost, redis.Port(testPort))
	redisLock, err := New(redisPool)
	if err != nil {
		panic(err)
	}
	//go testGoroutine(key, 1, number, redisLock)
	for i := 0; i < threadNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < number; i++ {
				ei := ParametersInterface{Params: i}
				err := redisLock.Lock(key, &ei)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "start")
	go func() {
		fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "count:", count)
		time.Sleep(10 * time.Second)
	}()
	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"), "App elapsed: ", elapsed, " --- ", expected, count)

}
