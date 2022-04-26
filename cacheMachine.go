package main

import (
	"fmt"
	"sync"
)

//===========[CACHE/STATIC]====================================================================================================

var caches map[string]Cacher

//===========[INTERFACES]====================================================================================================

//Key defines types that can be used as keys in the cache
type Key interface {
	string | int | int64 | int32 | int16 | int8 | float32 | float64
}

type Cacher interface {
	Cache() *Cache[TKey, TValue]
}

//===========[STRUCTS]====================================================================================================

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

//===========[FUNCTIONALITY]====================================================================================================

//New initiates new cache. The two arguments define what type key and value the cache is going to hold
func New[TKey Key, TValue any](name string, k TKey, v TValue) Cache[TKey, TValue] {
	c := Cache[TKey, TValue]{
		data: make(map[TKey]TValue),
		mx:   sync.RWMutex{},
	}

	caches[name] = &c

	return c
}

func main() {
	fmt.Println("asdad")
}
