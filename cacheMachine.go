package cacheMachine

import (
	"sync"
	"time"
)

//===========[CACHE/STATIC]=============================================================================================

//===========[INTERFACES]===============================================================================================

//Key defines types that can be used as keys in the cache
type Key interface {
	string | int | int64 | int32 | int16 | int8 | float32 | float64 | bool
}

type AllGetter[TKey Key, TValue any] interface {
	GetAll() map[TKey]TValue
}

type AllGetterAndRemover[TKey Key, TValue any] interface {
	GetAllAndRemove() map[TKey]TValue
}

type BulkAdder[TKey Key, TValue any] interface {
	AddBulk(d map[TKey]TValue)
}

//===========[STRUCTS]==================================================================================================

//Individual entry in the cache
type entry[TValue any] struct {
	//The value stored in the cache
	Value TValue `json:"value" bson:"value"`

	//When was the value added to the cache
	TimeAdded time.Time `json:"time_added" bson:"time_added"`

	//How long will the cache wait after the last update until it will remove this element
	TimeoutDuration time.Duration `json:"valid_until" bson:"valid_until"`

	//This is the timer that monitors auto-removal of the element
	timer *time.Timer
}

//ResetTimer resets the timer for auto removal of this element from the cache
func (e entry[TValue]) ResetTimer() {
	e.timer.Reset(e.TimeoutDuration)
}

//TODO: Add json encoding
//TODO: Instead of using TValue, create custom type "entry" and have it contain the Value, added date, auto removal, etc..

//Cache is the main definition of the cache
type Cache[TKey Key, TValue any] struct {
	defaultTimeoutDuration time.Duration
	data                   map[TKey]entry[TValue]
	mx                     sync.RWMutex
}

//------PRIVATE------

//add method adds an item. This method has no mutex protection
func (c *Cache[TKey, TValue]) add(key TKey, val TValue) {
	c.data[key] = entry[TValue]{
		Value:           val,
		TimeAdded:       time.Now(),
		TimeoutDuration: c.defaultTimeoutDuration,
		timer: time.AfterFunc(c.defaultTimeoutDuration, func() {
			//TODO: Check if this working correctly
			c.Remove(key)
		}),
	}
}

//remove method removes an item, but is not protected by a mutex
func (c *Cache[TKey, TValue]) remove(key TKey) {
	//If data doesn't exist, there's no need to perform further operations
	if _, exist := c.data[key]; !exist {
		return
	}

	delete(c.data, key)
}

//Creates a copy of the data. This function is not protected by locks
func (c *Cache[TKey, TValue]) copyData() map[TKey]TValue {
	cpy := make(map[TKey]TValue)
	for key, entry := range c.data {
		cpy[key] = entry.Value
	}
	return cpy
}

//reset clears the cache, but it's not using locks
func (c *Cache[TKey, TValue]) reset() {
	c.data = make(map[TKey]entry[TValue])
}

//------PUBLIC------

//Add inserts new Value into the cache
func (c Cache[TKey, TValue]) Add(key TKey, val TValue) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.add(key, val)
}

//AddBulk adds items to cache in bulk
func (c Cache[TKey, TValue]) AddBulk(d map[TKey]TValue) {
	if d == nil {
		return
	}

	c.mx.Lock()
	defer c.mx.Unlock()
	for k, v := range d {
		c.add(k, v)
	}
}

//Remove removes Value from the cache based on the key provided
func (c Cache[TKey, TValue]) Remove(key TKey) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.remove(key)
}

//RemoveBulk removes cached data based on keys provided
func (c Cache[TKey, TValue]) RemoveBulk(keys []TKey) {
	if keys == nil || len(keys) < 1 {
		return
	}

	c.mx.Lock()
	defer c.mx.Unlock()
	for _, key := range keys {
		c.remove(key)
	}
}

//Get returns Value based on the key provided
func (c Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	entry, exist := c.data[key]
	return entry.Value, exist
}

//GetBulk returns a map of key -> Value pairs where key is one provided in the slice
func (c Cache[TKey, TValue]) GetBulk(d []TKey) map[TKey]TValue {
	results := make(map[TKey]TValue)

	c.mx.RLock()
	for _, k := range d {
		results[k] = c.data[k].Value
	}
	c.mx.RUnlock()

	return results
}

//GetAndRemove returns requested Value and removes it from the cache
func (c Cache[TKey, TValue]) GetAndRemove(key TKey) (TValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	entry, exist := c.data[key]
	c.remove(key)
	return entry.Value, exist
}

//GetAll returns all the values stored in the cache
func (c Cache[TKey, TValue]) GetAll() map[TKey]TValue {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.copyData()
}

//GetAllAndRemove returns and removes all the elements from the cache
func (c Cache[TKey, TValue]) GetAllAndRemove() map[TKey]TValue {
	c.mx.Lock()
	defer c.mx.Unlock()
	cpy := c.copyData()
	c.reset()
	return cpy
}

//GetRandomSamples returns mixed set of items. Number of items is defined in the argument, if it exceeds the
//number of items that are present in the cache, it will return all the cached items
func (c Cache[TKey, TValue]) GetRandomSamples(n int) map[TKey]TValue {
	results := make(map[TKey]TValue)

	for key, entry := range c.data {
		if n < 1 {
			break
		}

		results[key] = entry.Value

		n--
	}

	return results
}

//Exist checks whether there the key exists in the cache
func (c Cache[TKey, TValue]) Exist(key TKey) bool {
	c.mx.RLock()
	defer c.mx.RUnlock()
	_, exist := c.data[key]
	return exist
}

//Count returns number of elements currently present in the cache
func (c Cache[TKey, TValue]) Count() int {
	c.mx.Lock()
	defer c.mx.Unlock()
	return len(c.data)
}

//ForEach runs a loop for each element in the cache. Take care using this method as it locks reading/writing the
//cache until ForEach completes.
func (c Cache[TKey, TValue]) ForEach(f func(TKey, TValue)) {
	d := c.GetAll()

	for k, v := range d {
		f(k, v)
	}
}

//Reset empties the cache and resets all the counters
func (c Cache[TKey, TValue]) Reset() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.reset()
}

//===========[FUNCTIONALITY]====================================================================================================

//New initiates new cache. It can also take in values that will be added to the cache immediately after initiation
func New[TKey Key, TValue any](initialValues map[TKey]TValue, timeout *time.Duration) Cache[TKey, TValue] {
	c := Cache[TKey, TValue]{
		data: make(map[TKey]entry[TValue]),
		mx:   sync.RWMutex{},
		defaultTimeoutDuration: timeout,
	}

	c.AddBulk(initialValues)

	return c
}

//Copy creates identical copy of the cache supplied as an argument
func Copy[TKey Key, TValue any](d AllGetter[TKey, TValue]) Cache[TKey, TValue] {
	return New[TKey, TValue](d.GetAll())
}

//Merge copies all data from cache2 into cache1
func Merge[TKey Key, TValue any](cache1 BulkAdder[TKey, TValue], cache2 AllGetter[TKey, TValue]) {
	cache1.AddBulk(cache2.GetAll())
}

//MergeAndReset copies all data from cache2 into cache1 and wipes cache2 clean right after
func MergeAndReset[TKey Key, TValue any](cache1 BulkAdder[TKey, TValue], cache2 AllGetterAndRemover[TKey, TValue]) {
	cache1.AddBulk(cache2.GetAllAndRemove())
}
