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
	time.Sleep(time.Second * 3)
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
	mmap.Start()                                       // try to Start when it has been already started
	mmap.Stop()                                        // stop
	log.Printf("length of mtxmap is %d\n", mmap.Len()) // length of mtxmap
	time.Sleep(time.Second * 10)                       // wait to ensure that mutexes expire
	log.Printf("length of mtxmap is %d\n", mmap.Len()) // still same length because cleaner has been stopped
	log.Printf("current state of db: %v", db)          // check global db variable
	mmap.Start()                                       // start cleaner again
	time.Sleep(time.Second * 2)                        // wait to that cleaner deletes expired mutexes
	log.Printf("length of mtxmap is %d\n", mmap.Len()) // see all mutexes has been deleted by cleaner

}

/*
2022/03/16 15:38:20 thread_id(14) acquire mutex of key2 key
2022/03/16 15:38:20 thread_id(3) acquire mutex of key1 key
2022/03/16 15:38:23 thread_id(3) releasing mutex of key1 key
2022/03/16 15:38:23 thread_id(14) releasing mutex of key2 key
2022/03/16 15:38:23 thread_id(4) acquire mutex of key1 key
2022/03/16 15:38:23 thread_id(11) acquire mutex of key2 key
2022/03/16 15:38:26 thread_id(4) releasing mutex of key1 key
2022/03/16 15:38:26 thread_id(11) releasing mutex of key2 key
2022/03/16 15:38:26 thread_id(1) acquire mutex of key1 key
2022/03/16 15:38:26 thread_id(10) acquire mutex of key2 key
2022/03/16 15:38:29 thread_id(1) releasing mutex of key1 key
2022/03/16 15:38:29 thread_id(2) acquire mutex of key1 key
2022/03/16 15:38:29 thread_id(10) releasing mutex of key2 key
2022/03/16 15:38:29 thread_id(13) acquire mutex of key2 key
2022/03/16 15:38:32 thread_id(2) releasing mutex of key1 key
2022/03/16 15:38:32 thread_id(13) releasing mutex of key2 key
2022/03/16 15:38:32 thread_id(12) acquire mutex of key2 key
2022/03/16 15:38:35 thread_id(12) releasing mutex of key2 key
2022/03/16 15:38:35 <mtxmap> mutex map already started!
2022/03/16 15:38:35 length of mtxmap is 2
2022/03/16 15:38:35 <mtxmap> map cleaner cancelled
2022/03/16 15:38:45 length of mtxmap is 2
2022/03/16 15:38:45 current state of db: map[key1:4 key2:5]
2022/03/16 15:38:46 <mtxmap> 2 of keys has been deleted
2022/03/16 15:38:47 length of mtxmap is 0
*/
