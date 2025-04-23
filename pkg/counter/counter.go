package counter

import (
	"sync"
	"sync/atomic"
)

var globalCounter safeCounter

type safeCounter struct{ counts sync.Map }

func Incr(key string, delta int64) {
	val, _ := globalCounter.counts.LoadOrStore(key, new(int64))
	atomic.AddInt64(val.(*int64), delta)
}

func Get(key string) int64 {
	val, ok := globalCounter.counts.Load(key)
	if !ok {
		return 0
	}
	return atomic.LoadInt64(val.(*int64))
}

func Reset() {
	globalCounter.counts.Range(func(key, _ interface{}) bool {
		globalCounter.counts.Delete(key)
		return true
	})
}
