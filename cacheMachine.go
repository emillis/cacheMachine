package cacheMachine

import (
	"sync"
)

//===========[CACHE/STATIC]=============================================================================================

//===========[INTERFACES]===============================================================================================

//Key defines types that can be used as keys in the cache
type Key interface {
	string | int | int64 | int32 | int16 | int8 | float32 | float64 | bool
}

//===========[STRUCTS]==================================================================================================

//Cache is the main definition of the cache
type Cache[TKey Key, TValue any] struct {
	data    map[TKey]TValue
	mx      sync.RWMutex
	counter int
}

//------PRIVATE------

//add method adds an item. This method has no mutex protection
func (c *Cache[TKey, TValue]) add(key TKey, val TValue) {
	if _, exist := c.data[key]; !exist {
		c.counter++
	}

	c.data[key] = val
}

//remove method removes an item, but is not protected by a mutex
func (c *Cache[TKey, TValue]) remove(key TKey) {
	//If data doesn't exist, there's no need to perform further operations
	if _, exist := c.data[key]; !exist {
		return
	}

	delete(c.data, key)

	c.counter--
}

//Creates a copy of the data
func (c *Cache[TKey, TValue]) copyData() map[TKey]TValue {
	cpy := make(map[TKey]TValue)
	c.mx.RLock()
	for key, val := range c.data {
		cpy[key] = val
	}
	c.mx.RUnlock()
	return cpy
}

//------PUBLIC------

//Add inserts new value into the cache
func (c *Cache[TKey, TValue]) Add(key TKey, val TValue) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.add(key, val)
}

//AddBulk adds items to cache in bulk
func (c *Cache[TKey, TValue]) AddBulk(d map[TKey]TValue) {
	if d == nil {
		return
	}

	c.mx.Lock()
	defer c.mx.Unlock()
	for k, v := range d {
		c.add(k, v)
	}
}

//Remove removes value from the cache based on the key provided
func (c *Cache[TKey, TValue]) Remove(key TKey) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.remove(key)
}

//RemoveBulk removes cached data based on keys provided
func (c *Cache[TKey, TValue]) RemoveBulk(keys []TKey) {
	if keys == nil || len(keys) < 1 {
		return
	}

	c.mx.Lock()
	defer c.mx.Unlock()
	for _, key := range keys {
		c.remove(key)
	}
}

//Get returns value based on the key provided
func (c *Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	v, exist := c.data[key]
	return v, exist
}

//GetBulk returns a map of key -> value pairs where key is one provided in the slice
func (c *Cache[TKey, TValue]) GetBulk(d []TKey) map[TKey]TValue {
	results := make(map[TKey]TValue)

	c.mx.RLock()
	for _, k := range d {
		results[k] = c.data[k]
	}
	c.mx.RUnlock()

	return results
}

//GetAndRemove returns requested value and removes it from the cache
func (c *Cache[TKey, TValue]) GetAndRemove(key TKey) (TValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	v, exist := c.data[key]
	c.remove(key)
	return v, exist
}

//GetAll returns all the values stored in the cache
func (c *Cache[TKey, TValue]) GetAll() map[TKey]TValue {
	return c.copyData()
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

//New initiates new cache. It can also take in values that will be added to the cache immediately after initiation
func New[TKey Key, TValue any](initialValues map[TKey]TValue) Cache[TKey, TValue] {
	c := Cache[TKey, TValue]{
		data: make(map[TKey]TValue),
		mx:   sync.RWMutex{},
	}

	c.AddBulk(initialValues)

	return c
}

//Copy creates identical copy of the cache supplied as an argument
func Copy[TKey Key, TValue any](c Cache[TKey, TValue]) Cache[TKey, TValue] {
	return New[TKey, TValue](c.GetAll())
}
