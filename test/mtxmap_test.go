package test

import (
	"github.com/bakyazi/mtxmap"
	"log"
	"sync"
	"testing"
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
	log.Printf("%d acquire %s key\n", id, key)
	db[key] += 1
	time.Sleep(time.Second * 1)
	log.Printf("%d releasing %s key\n", id, key)

}

func TestMtxmap(t *testing.T) {
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

	if mmap.Len() != 2 {
		t.Logf("failed map size. expected: %d output:%d\n", 2, mmap.Len())
		t.Fail()
		return
	}

	time.Sleep(time.Second * 10)

	if val, ok := db["key1"]; !ok || val != 4 {
		t.Logf("failed db key1 value. expected: %d output:%d\n", 4, val)
		t.Fail()
		return
	}
	if val, ok := db["key2"]; !ok || val != 4 {
		t.Logf("failed db key2 value. expected: %d output:%d\n", 4, val)
		t.Fail()
		return
	}
	if mmap.Len() != 0 {
		t.Logf("failed map size. expected: %d output:%d\n", 0, mmap.Len())
		t.Fail()
		return
	}

}
