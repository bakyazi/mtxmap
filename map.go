package mtxmap

import (
	"sync"
	"time"
)

// MutexMap is a map containing sync.Map whose value is mutexEntity
// it use as key based mutex map. Check if mutexEntities is expired every second
// if it is expired then deletes it from map
type MutexMap struct {
	TTL  time.Duration
	data sync.Map
	size int
}

// NewMutexMap creates and return new MutexMap object
func NewMutexMap(ttl time.Duration) *MutexMap {
	mmap := &MutexMap{
		TTL: ttl,
	}

	go func() {
		for range time.Tick(time.Second) {
			mmap.data.Range(func(key, value interface{}) bool {
				if value.(*mutexEntity).isExpire(mmap.TTL) {
					mmap.data.Delete(key)
					mmap.decrement()
					//log.Printf("deleting %v key\n", key)
				}
				return true
			})
		}
	}()
	return mmap
}

// Lock firstly creates a mutex if there is no mutex with given key
// if mutex with given key is already stored in map then retrieves it and
// tries to lock that mutex
//
// if a mutex is newlyu created then increments size
//
// return unlock function of the mutex of that key
func (m *MutexMap) Lock(key interface{}) func() {
	val, loaded := m.data.LoadOrStore(key, &mutexEntity{})
	mtx := val.(*mutexEntity)
	mtx.lock()
	if !loaded {
		m.increment()
	}
	return func() {
		mtx.unlock()
	}
}

// Unlock if there is a mutex with given key, unlocks it
func (m *MutexMap) Unlock(key interface{}) {
	val, ok := m.data.Load(key)
	if ok {
		val.(*mutexEntity).unlock()
	}
}

func (m *MutexMap) increment() {
	m.size++
}

func (m *MutexMap) decrement() {
	m.size--
}

func (m *MutexMap) Len() int {
	return m.size
}
