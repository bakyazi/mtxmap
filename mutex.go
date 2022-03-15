package mtxmap

import (
	"sync"
	"sync/atomic"
	"time"
)

// mutexEntity a custom sync.Mutex has attributes lastAccess and count
// lastAccess is last time when that mutex is used
// count total number of goroutines waiting/using that mutex
type mutexEntity struct {
	sync.Mutex
	lastAccess time.Time
	count      uint64
}

// lock increments count and locks mutex
func (m *mutexEntity) lock() {
	atomic.AddUint64(&m.count, 1)
	m.Lock()
}

// unlock set lastAccess, unlocks mutex, and decrement count
func (m *mutexEntity) unlock() {
	m.lastAccess = time.Now()
	m.Unlock()
	atomic.AddUint64(&m.count, ^uint64(0))
}

// isExpire if count is 0 (zero) and lastAccess + ttl is before time.Now returns true
// otherwise return false
func (m *mutexEntity) isExpire(ttl time.Duration) bool {
	return atomic.CompareAndSwapUint64(&m.count, 0, 0) &&
		m.lastAccess.Add(ttl).Before(time.Now())
}
