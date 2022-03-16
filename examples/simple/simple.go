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
	mmap = mtxmap.NewMutexMap(time.Second*5, time.Second)
	wg := sync.WaitGroup{}
	for i := 1; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			incrementKey("key1", x)
			wg.Done()
		}(i)
	}

	for i := 10; i < 15; i++ {
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
	log.Printf("current state of db: %v", db)
}

/*
2022/03/16 15:41:49 thread_id(14) acquire mutex of key2 key
2022/03/16 15:41:49 thread_id(3) acquire mutex of key1 key
2022/03/16 15:41:50 thread_id(14) releasing mutex of key2 key
2022/03/16 15:41:50 thread_id(3) releasing mutex of key1 key
2022/03/16 15:41:50 thread_id(4) acquire mutex of key1 key
2022/03/16 15:41:50 thread_id(13) acquire mutex of key2 key
2022/03/16 15:41:51 thread_id(4) releasing mutex of key1 key
2022/03/16 15:41:51 thread_id(1) acquire mutex of key1 key
2022/03/16 15:41:51 thread_id(13) releasing mutex of key2 key
2022/03/16 15:41:51 thread_id(10) acquire mutex of key2 key
2022/03/16 15:41:52 thread_id(10) releasing mutex of key2 key
2022/03/16 15:41:52 thread_id(11) acquire mutex of key2 key
2022/03/16 15:41:52 thread_id(1) releasing mutex of key1 key
2022/03/16 15:41:52 thread_id(2) acquire mutex of key1 key
2022/03/16 15:41:53 thread_id(2) releasing mutex of key1 key
2022/03/16 15:41:53 thread_id(11) releasing mutex of key2 key
2022/03/16 15:41:53 thread_id(12) acquire mutex of key2 key
2022/03/16 15:41:54 thread_id(12) releasing mutex of key2 key
2022/03/16 15:41:54 length of mtxmap is 2
2022/03/16 15:41:59 <mtxmap> 1 of keys has been deleted
2022/03/16 15:42:00 <mtxmap> 1 of keys has been deleted
2022/03/16 15:42:04 length of mtxmap is 0
2022/03/16 15:42:04 current state of db: map[key1:4 key2:5]
*/
