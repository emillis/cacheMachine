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
	ResetTimer()
	StopTimer()
	TimerExist() bool
	TimeLeft() time.Time
	SetTimer(duration time.Duration)
	SetAndResetTimer(duration time.Duration)
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

	//When was the value added to the cache
	TimeAdded time.Time `json:"time_added" bson:"time_added"`

	//How long will the cache wait after the last update until it will remove this element
	TimeoutDuration time.Duration `json:"valid_until" bson:"valid_until"`

	//This is the timer that monitors auto-removal of the element
	timer *time.Timer

	//Stores time when the last timer reset happened
	lastReset time.Time

	//Locks
	mx sync.RWMutex
}

//------PRIVATE------

//Sets a new duration after which the entity is removed. This method is not protected by a mutex
func (e *entry[TValue]) newTimeoutDuration(duration time.Duration) {
	e.TimeoutDuration = duration
}

//Resets timeout duration back to the beginning. This is not protected by a mutex
func (e *entry[TValue]) resetTimeout() {
	if e.timer == nil {
		return
	}

	e.timer.Reset(e.TimeoutDuration)
	e.lastReset = time.Now()
}

//------PUBLIC------

//Value returns the value of this entry
func (e *entry[TValue]) Value() TValue {
	return e.Val
}

//ResetTimer resets the countdown timer until the removal of this entry
func (e *entry[TValue]) ResetTimer() {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.resetTimeout()
}

//TimerExist returns time left until removal of the entity
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
	e.timer.Stop()
}

//TimeLeft returns time when the entity will be removed
func (e *entry[TValue]) TimeLeft() time.Time {
	return e.lastReset.Add(e.TimeoutDuration)
}

//SetTimer sets new duration after which entity is removed. This method does not force entity to start using the new
//timeout duration. To achieve that, you need to either use "SetAndResetTimer" method or additionally call
//"ResetTimer" method
func (e *entry[TValue]) SetTimer(duration time.Duration) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.newTimeoutDuration(duration)
}

//SetAndResetTimer does exactly what SetTimer, but also resets the timeout
func (e *entry[TValue]) SetAndResetTimer(duration time.Duration) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.newTimeoutDuration(duration)
	e.resetTimeout()
}

//TODO: Add json encoding

//Cache is the main definition of the cache
type cache[TKey Key, TValue any] struct {
	Requirements Requirements
	data         map[TKey]entry[TValue]
	mx           sync.RWMutex
}
type Cache[TKey Key, TValue any] struct {
	cache[TKey, TValue]
}

//------PRIVATE------

//add method adds an item. This method has no mutex protection
func (c *Cache[TKey, TValue]) add(key TKey, val TValue) Entry[TValue] {
	now := time.Now()

	e := entry[TValue]{
		Val:             val,
		TimeAdded:       now,
		TimeoutDuration: c.cache.Requirements.DefaultTimeout,
		mx:              sync.RWMutex{},
	}

	if c.cache.Requirements.timeoutInUse {
		e.timer = time.AfterFunc(c.cache.Requirements.DefaultTimeout, func() {
			c.Remove(key)
		})

		e.lastReset = now
	}

	c.data[key] = e

	return &e
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
	c.data = make(map[TKey]entry[TValue])
}

//------PUBLIC------

//Add inserts new Val into the cache
func (c Cache[TKey, TValue]) Add(key TKey, val TValue) Entry[TValue] {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.add(key, val)
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
	e, exist := c.data[key]
	c.remove(key)
	return e.Val, exist
}

//GetAll returns all the values stored in the cache
func (c Cache[TKey, TValue]) GetAll() map[TKey]TValue {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.copyValues()
}

//GetAllAndRemove returns and removes all the elements from the cache
func (c *Cache[TKey, TValue]) GetAllAndRemove() map[TKey]TValue {
	c.mx.Lock()
	defer c.mx.Unlock()
	cpy := c.copyValues()
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
		data:         make(map[TKey]entry[TValue]),
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
