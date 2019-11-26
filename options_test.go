package distlock

import (
	"fmt"
	"testing"
	"time"
)

func TestBackOffDelay(t *testing.T) {
	configOption := &ConfigOption{delay: 200 * time.Millisecond}
	var n uint
	for n < 30 {
		fmt.Println(BackOffDelay(n, configOption))
		n++
	}
}

func TestGeneratorLockKey(t *testing.T) {
	test := GeneratorLockKey("poiuytrewq0")
	fmt.Println(test)
	test = GeneratorLockKey("poiuytrewq1")
	fmt.Println(test)
	test = GeneratorLockKey("poiuytrewq2")
	fmt.Println(test)
}
