package distlock

import (
	"encoding/base64"
	"github.com/rs/xid"
	"math/rand"
	"time"
)

func GeneratorRequestId() string {
	return xid.New().String()
}

func BackOffDelay(n uint, co *ConfigOption) time.Duration {
	return co.delay + time.Duration(rand.Intn(100))*time.Millisecond
	//result := co.delay * (1 << n)
	//if result > 2*time.Second {
	//	return time.Second
	//}
	//return result
	//return co.delay * (1 << n)
}

func GeneratorLockKey(lockKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(lockKey))
}
