# mtxmap
key based mutex map with ttl

## example

```go
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

func useMap(key string, id int) {
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
			useMap("key1", x)
			wg.Done()
		}(i)
	}

	for i := 1; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			useMap("key2", x)
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

```

output

```
2022/03/15 21:02:48 thread_id(4) acquire mutex of key2 key
2022/03/15 21:02:48 thread_id(3) acquire mutex of key1 key
2022/03/15 21:02:49 thread_id(4) releasing mutex of key2 key
2022/03/15 21:02:49 thread_id(2) acquire mutex of key2 key
2022/03/15 21:02:49 thread_id(3) releasing mutex of key1 key
2022/03/15 21:02:49 thread_id(1) acquire mutex of key1 key
2022/03/15 21:02:50 thread_id(2) releasing mutex of key2 key
2022/03/15 21:02:50 thread_id(3) acquire mutex of key2 key
2022/03/15 21:02:50 thread_id(1) releasing mutex of key1 key
2022/03/15 21:02:50 thread_id(2) acquire mutex of key1 key
2022/03/15 21:02:51 thread_id(3) releasing mutex of key2 key
2022/03/15 21:02:51 thread_id(1) acquire mutex of key2 key
2022/03/15 21:02:51 thread_id(2) releasing mutex of key1 key
2022/03/15 21:02:51 thread_id(4) acquire mutex of key1 key
2022/03/15 21:02:52 thread_id(1) releasing mutex of key2 key
2022/03/15 21:02:52 thread_id(4) releasing mutex of key1 key
2022/03/15 21:02:52 length of mtxmap is 2
2022/03/15 21:03:02 length of mtxmap is 0
2022/03/15 21:03:02 value of key key1 is 4
2022/03/15 21:03:02 value of key key2 is 4

```