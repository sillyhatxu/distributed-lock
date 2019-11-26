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
	test := GeneratorLockKey("C01D7CF6EC3F47F09556A5D6E9009A43")
	fmt.Println(test)
	test = GeneratorLockKey("poiuytrewq1")
	fmt.Println(test)
	test = GeneratorLockKey("poiuytrewq2")
	fmt.Println(test)
}
