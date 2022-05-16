package main

import (
	"fmt"
	"sync"
	"time"
)

//===========[CACHE/STATIC]=============================================================================================

//===========[INTERFACES]===============================================================================================

//Key defines types that can be used as keys in the cache
type Key interface {
	string | int | int64 | int32 | int16 | int8 | float32 | float64
}

type SaveHandler[TKey Key, TValue any] interface {
	Save(map[TKey]TValue)
}

//===========[STRUCTS]==================================================================================================

//Cache is the main definition of the cache
type Cache[TKey Key, TValue any] struct {
	data    map[TKey]TValue
	mx      sync.RWMutex
	counter int
}

//Add inserts new value into the cache
func (c *Cache[TKey, TValue]) Add(key TKey, val TValue) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data[key] = val
	c.counter++
}

//Remove removes value from the cache based on the key provided
func (c *Cache[TKey, TValue]) Remove(key TKey) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.data, key)
	c.counter--
}

//Get returns value based on the key provided
func (c *Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	v, exist := c.data[key]
	return v, exist
}

//Exist checks whether there the key exists in the cache
func (c *Cache[TKey, TValue]) Exist(key TKey) bool {
	c.mx.RLock()
	defer c.mx.RUnlock()
	_, exist := c.data[key]
	return exist
}

//Count returns number of elements currently present in the cache
func (c *Cache[TKey, TValue]) Count() int {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.counter
}

//Reset empties the cache and resets all the counters
func (c *Cache[TKey, TValue]) Reset() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data = make(map[TKey]TValue)
	c.counter = 0
}

//Creates a copy of the data
func (c *Cache[TKey, TValue]) copyData() map[TKey]TValue {
	cpy := make(map[TKey]TValue)
	c.mx.Lock()
	for key, val := range c.data {
		cpy[key] = val
	}
	c.mx.Unlock()
	return cpy
}

//Save saves the entire cache
func (c *Cache[TKey, TValue]) Save(saveHandler SaveHandler[TKey, TValue]) {
	saveHandler.Save(c.copyData())
}

func (c *Cache[TKey, TValue]) SetSaveInterval(saveHandler SaveHandler[TKey, TValue], ticker *time.Ticker) {
	go func() {
		for {
			<-ticker.C
		}
	}()
}

//===========[FUNCTIONALITY]====================================================================================================

//New initiates new cache. The two arguments define what type key and value the cache is going to hold
func New[TKey Key, TValue any](name string, k TKey, v TValue) Cache[TKey, TValue] {
	c := Cache[TKey, TValue]{
		data: make(map[TKey]TValue),
		mx:   sync.RWMutex{},
	}

	return c
}

func main() {
	cache := New("test_cache", "", 1)

	cache.Add("one", 77)
	cache.Add("two", 5)

	fmt.Println(cache.Count())
	fmt.Println(cache.Get("two"))

	cache.SetSaveInterval(nil, time.NewTicker(time.Second * 2))
}
