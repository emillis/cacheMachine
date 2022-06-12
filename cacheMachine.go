package cacheMachine

//===========[INTERFACES]===============================================================================================

//Key defines types that can be used as keys in the Cache
type Key interface {
	string | int | int64 | int32 | int16 | int8 | float32 | float64
}

//===========[STRUCTS]==================================================================================================

type Cache[TKey Key, TValue any] struct {
	data             map[TKey]TValue
	counter          int
	addChan          chan map[TKey]TValue
	removeChan       chan []TKey
	resetChan        chan struct{}
	existChan        chan TKey
	returnExistChan  chan bool
	getCountChan     chan struct{}
	returnCountChan  chan int
	getAllChan       chan struct{}
	returnGetAllChan chan map[TKey]TValue
	getBulkChan      chan []TKey
	returnBulkCHan   chan map[TKey]TValue
	getSingleChan    chan TKey
	returnSingleChan chan TValue
}

//PRIVATE
//PRIVATE
//PRIVATE

//Add inserts new value into the Cache
func (c *Cache[TKey, TValue]) add(key TKey, val TValue) {
	if _, exist := c.data[key]; !exist {
		c.counter++
	}

	c.data[key] = val
}

//AddBulk adds items to Cache in bulk
func (c *Cache[TKey, TValue]) addBulk(d map[TKey]TValue) {
	if d == nil {
		return
	}

	for k, v := range d {
		c.add(k, v)
	}
}

//Remove removes value from the Cache based on the key provided
func (c *Cache[TKey, TValue]) remove(key TKey) {
	//If data doesn't exist, there's no need to perform further operations
	if _, exist := c.data[key]; !exist {
		return
	}

	delete(c.data, key)

	c.counter--
}

//RemoveBulk removes cached data based on keys provided
func (c *Cache[TKey, TValue]) removeBulk(keys []TKey) {
	if keys == nil || len(keys) < 1 {
		return
	}

	for _, key := range keys {
		//If data doesn't exist, there's no need to perform further commands
		c.remove(key)
	}
}

//Get returns value based on the key provided
func (c *Cache[TKey, TValue]) get(key TKey) TValue {
	return c.data[key]
}

//GetBulk returns a map of key -> value pairs where key is one provided in the slice
func (c *Cache[TKey, TValue]) getBulk(d []TKey) map[TKey]TValue {
	results := make(map[TKey]TValue)

	if d == nil || len(d) < 1 {
		return results
	}

	for _, k := range d {
		results[k] = c.data[k]
	}

	return results
}

//GetAll returns all the values stored in the Cache
func (c *Cache[TKey, TValue]) getAll() map[TKey]TValue {
	results := make(map[TKey]TValue)

	for k, v := range c.data {
		results[k] = v
	}

	return results
}

//Exist checks whether there the key exists in the Cache
func (c *Cache[TKey, TValue]) exist(key TKey) bool {
	_, exist := c.data[key]
	return exist
}

//Count returns number of elements currently present in the Cache
func (c *Cache[TKey, TValue]) count() int {
	return c.counter
}

//Reset empties the Cache and resets all the counters
func (c *Cache[TKey, TValue]) reset() {
	c.data = make(map[TKey]TValue)
	c.counter = 0
}

//PUBLIC
//PUBLIC
//PUBLIC

//Add inserts new value into the Cache
func (c *Cache[TKey, TValue]) Add(key TKey, val TValue) {
	c.addChan <- map[TKey]TValue{key: val}
}

//AddBulk adds items to Cache in bulk
func (c *Cache[TKey, TValue]) AddBulk(d map[TKey]TValue) {
	c.addChan <- d
}

//Remove removes value from the Cache based on the key provided
func (c *Cache[TKey, TValue]) Remove(key TKey) {
	c.removeChan <- []TKey{key}
}

//RemoveBulk removes cached data based on keys provided
func (c *Cache[TKey, TValue]) RemoveBulk(keys []TKey) {
	c.removeChan <- keys
}

//Get returns value based on the key provided
func (c *Cache[TKey, TValue]) Get(key TKey) TValue {
	c.getSingleChan <- key
	return <-c.returnSingleChan
}

//GetBulk returns a map of key -> value pairs where key is one provided in the slice
func (c *Cache[TKey, TValue]) GetBulk(d []TKey) map[TKey]TValue {
	c.getBulkChan <- d
	return <-c.returnBulkCHan
}

//GetAll returns all the values stored in the Cache
func (c *Cache[TKey, TValue]) GetAll() map[TKey]TValue {
	c.getAllChan <- struct{}{}
	return <-c.returnGetAllChan
}

//Exist checks whether there the key exists in the Cache
func (c *Cache[TKey, TValue]) Exist(key TKey) bool {
	c.existChan <- key
	return <-c.returnExistChan
}

//Count returns number of elements currently present in the Cache
func (c *Cache[TKey, TValue]) Count() int {
	c.getCountChan <- struct{}{}
	return <-c.returnCountChan
}

//Reset empties the Cache and resets all the counters
func (c *Cache[TKey, TValue]) Reset() {
	c.resetChan <- struct{}{}
}

//===========[FUNCTIONALITY]============================================================================================

//cacheManager is spawned as a goroutine for each new Cache
func cacheManager[TKey Key, TValue any](c *Cache[TKey, TValue]) {
	if c == nil {
		return
	}

	for {
		select {
		case addMap := <-c.addChan:
			c.addBulk(addMap)

		case removeSlice := <-c.removeChan:
			c.removeBulk(removeSlice)

		case <-c.resetChan:
			c.reset()

		case key := <-c.existChan:
			c.returnExistChan <- c.exist(key)

		case <-c.getCountChan:
			c.returnCountChan <- c.count()

		case <-c.getAllChan:
			c.returnGetAllChan <- c.getAll()

		case keys := <-c.getBulkChan:
			c.returnBulkCHan <- c.getBulk(keys)

		case key := <-c.getSingleChan:
			c.returnSingleChan <- c.get(key)

		}
	}
}

//New initiates new Cache. The two arguments define what type key and value the Cache is going to hold
func New[TKey Key, TValue any]() Cache[TKey, TValue] {
	c := Cache[TKey, TValue]{
		data:             make(map[TKey]TValue),
		counter:          0,
		addChan:          make(chan map[TKey]TValue),
		removeChan:       make(chan []TKey),
		resetChan:        make(chan struct{}),
		existChan:        make(chan TKey),
		returnExistChan:  make(chan bool),
		getCountChan:     make(chan struct{}),
		returnCountChan:  make(chan int),
		getAllChan:       make(chan struct{}),
		returnGetAllChan: make(chan map[TKey]TValue),
		getBulkChan:      make(chan []TKey),
		returnBulkCHan:   make(chan map[TKey]TValue),
		getSingleChan:    make(chan TKey),
		returnSingleChan: make(chan TValue),
	}

	go cacheManager[TKey, TValue](&c)

	return c
}
