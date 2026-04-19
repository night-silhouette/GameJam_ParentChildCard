package Util

import "time"

func SendTime(offset time.Duration) int64 {
	return time.Now().Add(offset).UnixMilli()
}

// CreateTimer 返回参数:(StopChan,CrashChan)
func CreateTimer(duration time.Duration, callback func()) (chan struct{}, chan struct{}) {
	timer := time.NewTimer(duration)
	StopChan := make(chan struct{}, 1)
	CrashChan := make(chan struct{}, 1)
	go func() {
		defer timer.Stop()
		select {
		case <-timer.C:
			callback()
		case <-StopChan:
			callback()
		case <-CrashChan:
		}

	}()
	return StopChan, CrashChan
}
