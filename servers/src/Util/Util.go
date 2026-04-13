package Util

import "time"

func SendTime(offset time.Duration) int64 {
	return time.Now().Add(offset).UnixMilli()
}
