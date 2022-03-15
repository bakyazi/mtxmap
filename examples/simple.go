package main

import (
	"github.com/bakyazi/mtxmap"
	"log"
	"sync"
	"time"
)

var mmap *mtxmap.MutexMap
var db = map[string]int{
	"key1": 0,
	"key2": 0,
}

func incrementKey(key string, id int) {
	unlock := mmap.Lock(key)
	defer unlock()
	log.Printf("thread_id(%d) acquire mutex of %s key\n", id, key)
	db[key] += 1
	time.Sleep(time.Second * 1)
	log.Printf("thread_id(%d) releasing mutex of %s key\n", id, key)
}

func main() {
	mmap = mtxmap.NewMutexMap(time.Second * 5)
	wg := sync.WaitGroup{}
	for i := 1; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			incrementKey("key1", x)
			wg.Done()
		}(i)
	}

	for i := 1; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			incrementKey("key2", x)
			wg.Done()
		}(i)
	}

	wg.Wait()

	log.Printf("length of mtxmap is %d\n", mmap.Len())
	time.Sleep(time.Second * 10)
	log.Printf("length of mtxmap is %d\n", mmap.Len())
	log.Printf("value of key %s is %d", "key1", db["key1"])
	log.Printf("value of key %s is %d", "key2", db["key2"])
}
