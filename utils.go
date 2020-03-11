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
	return co.delay + (time.Duration(rand.Intn(int(n*10))) * time.Millisecond)
}

func GeneratorLockKey(lockKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(lockKey))
}
