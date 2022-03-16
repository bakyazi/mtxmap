package mtxmap

import (
	"context"
	"log"
	"reflect"
	"sync"
	"time"
)

// MutexMap is a map containing sync.Map whose value is mutexEntity
// it use as key based mutex map. Check if mutexEntities is expired every second
// if it is expired then deletes it from map
type MutexMap struct {
	ttl           time.Duration
	cleanInterval time.Duration
	size          int
	mutex         *sync.Mutex
	ctx           context.Context
	cancelFunc    context.CancelFunc
	data          sync.Map
}

// NewMutexMap creates and return new MutexMap object
func NewMutexMap(ttl time.Duration, cleanInterval time.Duration) *MutexMap {
	m := &MutexMap{
		mutex:         &sync.Mutex{},
		ttl:           ttl,
		cleanInterval: cleanInterval,
	}
	m.ctx, m.cancelFunc = context.WithCancel(context.Background())
	m.Start()

	return m
}

// Lock firstly creates a mutex if there is no mutex with given key
// if mutex with given key is already stored in map then retrieves it and
// tries to lock that mutex
//
// if a mutex is newly created then increments size
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

func (m *MutexMap) Len() int {
	return m.size
}

func (m *MutexMap) SetTTL(t time.Duration) {
	m.ttl = t
}

// Start starts cleaner goroutine
func (m *MutexMap) Start() {
	if isMutexLocked(m.mutex) {
		log.Println("<mtxmap> mutex map already started!")
		return
	}
	go func() {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.ctx, m.cancelFunc = context.WithCancel(context.Background())
		ticker := time.Tick(m.cleanInterval)
		for {
			select {
			case <-ticker:
				deletedCount := 0
				m.data.Range(func(key, value interface{}) bool {
					if value.(*mutexEntity).isExpire(m.ttl) {
						m.data.Delete(key)
						m.decrement()
						deletedCount++
					}
					return true
				})
				if deletedCount > 0 {
					log.Printf("<mtxmap> %d of keys has been deleted\n", deletedCount)
				}
			case <-m.ctx.Done():
				log.Println("<mtxmap> map cleaner cancelled")
				return
			}
		}
	}()
}

func (m *MutexMap) Stop() {
	if m.cancelFunc == nil {
		return
	}
	m.cancelFunc()
	m.ctx = nil
	m.cancelFunc = nil
}

func (m *MutexMap) increment() {
	m.size++
}

func (m *MutexMap) decrement() {
	m.size--
}

func isMutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&1 == 1
}
