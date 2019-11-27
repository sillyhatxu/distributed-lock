package distlock

import "errors"

var ErrFailed = errors.New("redis lock: failed to acquire lock")
var ErrTimeOut = errors.New("redis lock: time out to acquire lock")
