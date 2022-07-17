package cacheMachine

import (
	"sync"
	"time"
)

//===========[CACHE/STATIC]=============================================================================================

var defaultRequirements = Requirements{}

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

type Entry[TValue any] interface {
	Value() TValue
	ResetTimer(time.Duration)
	StopTimer()
	TimerExist() bool
}

//===========[STRUCTS]==================================================================================================

type Requirements struct {
	//If this is set, by default, every cache entry will have a timeout of this duration after which
	//the element will be removed from the cache. This timeout can be changed for individual entry
	DefaultTimeout time.Duration

	//Defines whether the DefaultTimeout is in use
	timeoutInUse bool
}

//Individual entry in the cache
type entry[TValue any] struct {
	//The value stored in the cache
	Val TValue `json:"value" bson:"value"`

	//This is the timer that monitors auto-removal of the element
	timer *time.Timer

	//Locks
	mx sync.RWMutex
}

//------PRIVATE------

//Resets timeout duration to the duration specified. If 0 is supplied, it stops the timer
func (e *entry[TValue]) resetTimer(t time.Duration) {
	if e.timer == nil {
		return
	}

	if t.String() == "0s" {
		e.timer.Stop()
		return
	}

	e.timer.Reset(t)
}

//------PUBLIC------

//Value returns the value of this entry
func (e *entry[TValue]) Value() TValue {
	return e.Val
}

//ResetTimer resets the countdown timer until the removal of this entry
func (e *entry[TValue]) ResetTimer(t time.Duration) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.resetTimer(t)
}

//TimerExist checks whether the timer exist and returns boolean accordingly
func (e *entry[TValue]) TimerExist() bool {
	if e.timer != nil {
		return true
	}

	return false
}

//StopTimer stops the countdown timer until the element is removed
func (e *entry[TValue]) StopTimer() {
	if e.timer == nil {
		return
	}

	e.mx.Lock()
	defer e.mx.Unlock()
	e.resetTimer(0)
}

//Cache is the main definition of the cache
type cache[TKey Key, TValue any] struct {
	Requirements Requirements
	data         map[TKey]*entry[TValue]
	mx           sync.RWMutex
}
type Cache[TKey Key, TValue any] struct {
	cache[TKey, TValue]
}

//------PRIVATE------

//add method adds an item. This method has no mutex protection
func (c *Cache[TKey, TValue]) add(key TKey, val TValue, t time.Duration) Entry[TValue] {
	e := entry[TValue]{
		Val: val,
		mx:  sync.RWMutex{},
	}

	//Timer implementation
	if t.String() != "0s" || c.cache.Requirements.timeoutInUse {
		if t.String() == "0s" {
			t = c.cache.Requirements.DefaultTimeout
		}

		e.timer = time.AfterFunc(t, func() {
			c.Remove(key)
		})
	}

	c.data[key] = &e

	return &e
}

//addTImer adds new timer with specified duration if it doesn't yet exist. If timer is already present,
//this method resets it with the specified duration
func (c *Cache[TKey, TValue]) addTimer(key TKey, t time.Duration) {
	e, exist := c.data[key]

	if !exist {
		return
	}

	if e.timer != nil {
		e.timer.Reset(t)
		return
	}

	e.timer = time.AfterFunc(t, func() { c.Remove(key) })
}

//remove method removes an item, but is not protected by a mutex
func (c *Cache[TKey, TValue]) remove(key TKey) {
	delete(c.data, key)
}

//Creates a copy of the data. This function is not protected by locks
func (c *Cache[TKey, TValue]) copyValues() map[TKey]TValue {
	cpy := make(map[TKey]TValue)
	for key, entry := range c.data {
		cpy[key] = entry.Val
	}
	return cpy
}

//reset clears the cache, but it's not using locks
func (c *Cache[TKey, TValue]) reset() {
	c.data = make(map[TKey]*entry[TValue])
}

//------PUBLIC------

//AddTimer adds timer to the key specified. If the key already has a timer, it gets reset with the new duration specified
func (c Cache[TKey, TValue]) AddTimer(key TKey, t time.Duration) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.addTimer(key, t)
}

//Add inserts new key:value pair into the cache
func (c Cache[TKey, TValue]) Add(key TKey, val TValue) Entry[TValue] {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.add(key, val, 0)
}

//AddWithTimeout does the same as method "Add" but also sets timer for automatic removal of the entry
func (c Cache[TKey, TValue]) AddWithTimeout(key TKey, val TValue, timeout time.Duration) Entry[TValue] {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.add(key, val, timeout)
}

//AddBulk adds items to cache in bulk
func (c Cache[TKey, TValue]) AddBulk(d map[TKey]TValue) {
	if d == nil {
		return
	}

	c.mx.Lock()
	defer c.mx.Unlock()
	for k, v := range d {
		c.add(k, v, 0)
	}
}

//Remove removes Val from the cache based on the key provided
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

//Get returns Val based on the key provided
func (c Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	entry, exist := c.data[key]
	return entry.Val, exist
}

//GetEntry returns Entry interface for the value saved in the cache
func (c Cache[TKey, TValue]) GetEntry(key TKey) Entry[TValue] {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.data[key]
}

//GetBulk returns a map of key -> Val pairs where key is one provided in the slice
func (c Cache[TKey, TValue]) GetBulk(d []TKey) map[TKey]TValue {
	results := make(map[TKey]TValue)

	c.mx.RLock()
	for _, k := range d {
		results[k] = c.data[k].Val
	}
	c.mx.RUnlock()

	return results
}

//GetAndRemove returns requested Val and removes it from the cache
func (c Cache[TKey, TValue]) GetAndRemove(key TKey) (TValue, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	defer c.remove(key)
	e, exist := c.data[key]
	return e.Val, exist
}

//GetAndRemoveEntry returns Entry interface and removes the entity from the cache immediately
func (c Cache[TKey, TValue]) GetAndRemoveEntry(key TKey) Entry[TValue] {
	c.mx.Lock()
	defer c.mx.Unlock()
	defer c.remove(key)
	return c.data[key]
}

//GetAll returns all the values stored in the cache
func (c Cache[TKey, TValue]) GetAll() map[TKey]TValue {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.copyValues()
}

//GetAllAndRemove returns and removes all the elements from the cache
func (c Cache[TKey, TValue]) GetAllAndRemove() map[TKey]TValue {
	c.mx.Lock()
	defer c.mx.Unlock()
	defer c.reset()
	return c.copyValues()
}

//GetRandomSamples returns mixed set of items. Number of items is defined in the argument, if it exceeds the
//number of items that are present in the cache, it will return all the cached items
func (c Cache[TKey, TValue]) GetRandomSamples(n int) map[TKey]TValue {
	results := make(map[TKey]TValue)

	for key, entry := range c.data {
		if n < 1 {
			break
		}

		results[key] = entry.Val

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
func (c *Cache[TKey, TValue]) Reset() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.reset()
}

//Requirements returns requirements used from this cache
func (c Cache[TKey, TValue]) Requirements() Requirements {
	return c.cache.Requirements
}

//===========[FUNCTIONALITY]====================================================================================================

//Adjusts and parses the Requirements
func makeRequirementsSensible(r *Requirements) {
	//Checking whether the DefaultTimeout is in use. If yes, it sets timeoutInUse to true
	r.timeoutInUse = r.DefaultTimeout.String() != "0s"
}

//New initiates new cache. It can also take in values that will be added to the cache immediately after initiation
func New[TKey Key, TValue any](r *Requirements) Cache[TKey, TValue] {
	if r == nil {
		r = &defaultRequirements
	}

	makeRequirementsSensible(r)

	c := cache[TKey, TValue]{
		Requirements: *r,
		data:         make(map[TKey]*entry[TValue]),
		mx:           sync.RWMutex{},
	}

	return Cache[TKey, TValue]{c}
}

//Copy creates identical copy of the cache supplied as an argument
func Copy[TKey Key, TValue any](c Cache[TKey, TValue]) Cache[TKey, TValue] {
	req := c.Requirements()
	nc := New[TKey, TValue](&req)
	nc.AddBulk(c.GetAll())
	return nc
}

//Merge copies all data from cache2 into cache1
func Merge[TKey Key, TValue any](cache1 BulkAdder[TKey, TValue], cache2 AllGetter[TKey, TValue]) {
	cache1.AddBulk(cache2.GetAll())
}

//MergeAndReset copies all data from cache2 into cache1 and wipes cache2 clean right after
func MergeAndReset[TKey Key, TValue any](cache1 BulkAdder[TKey, TValue], cache2 AllGetterAndRemover[TKey, TValue]) {
	cache1.AddBulk(cache2.GetAllAndRemove())
}
