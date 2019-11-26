package distlock

import (
	"encoding/base64"
	"time"
)

func BackOffDelay(n uint, co *ConfigOption) time.Duration {
	result := co.delay * (1 << n)
	if result > 2*time.Second {
		return time.Second
	}
	return result
	//return co.delay * (1 << n)
}

func GeneratorLockKey(lockKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(lockKey))
}
