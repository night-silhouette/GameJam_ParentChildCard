package Util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type RequestInterceptor struct {
	// key 是消息内容的 MD5, value 是收到该消息的时间
	history   map[string]time.Time
	mu        sync.Mutex
	threshold time.Duration
}

func NewInterceptor(duration time.Duration) *RequestInterceptor {
	return &RequestInterceptor{
		history:   make(map[string]time.Time),
		threshold: duration,
	}
}

func (l *RequestInterceptor) ShouldBlock(msg []byte) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	hash := md5.Sum(msg)
	fingerprint := hex.EncodeToString(hash[:])

	now := time.Now()

	for h, t := range l.history {
		if now.Sub(t) > l.threshold {
			delete(l.history, h)
		}
	}

	if lastTime, exists := l.history[fingerprint]; exists {
		if now.Sub(lastTime) < l.threshold {
			fmt.Printf("[拦截拦截器] 发现高频重复消息，已成功拦截！原始文本: %s\n", string(msg))
			return true
		}
	}

	l.history[fingerprint] = now
	return false
}
