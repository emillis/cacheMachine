package cacheMachine

import (
	"testing"
	"time"
)

//===========[FUNCTIONALITY]====================================================================================================

func initializeFullCache(n int, r *Requirements) Cache[int, int] {
	c := New[int, int](r)

	for i := 0; i < n; i++ {
		c.Add(i, i)
	}

	return c
}

//===========[TESTING]====================================================================================================

func TestCache_Add(t *testing.T) {
	c := initializeFullCache(10, nil)

	dataLength := len(c.data)

	if dataLength != 10 {
		t.Errorf("Expected value %d, received %d", 10, dataLength)
	}
}

func TestCache_AddBulk(t *testing.T) {
	c := initializeFullCache(0, nil)

	expectedLength := 5

	c.AddBulk(map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	})

	dataLength := len(c.data)

	if dataLength != expectedLength {
		t.Errorf("Expected value %d, received %d", expectedLength, dataLength)
	}
}

func TestCache_Count(t *testing.T) {
	expectedLength := 10

	c := initializeFullCache(expectedLength, nil)

	if c.Count() != expectedLength {
		t.Errorf("Expected value %d, received %d", expectedLength, c.Count())
	}
}

func TestCache_Get(t *testing.T) {
	requiredValue := 5

	c := initializeFullCache(10, nil)

	v, ok := c.Get(requiredValue)

	if v != requiredValue || !ok {
		t.Errorf("Required value was %d and %t, received %d and %t", requiredValue, true, v, ok)
	}
}

func TestCache_Exist(t *testing.T) {
	c := initializeFullCache(10, nil)

	requiredValue := 5

	if !c.Exist(requiredValue) {
		t.Errorf("Value %d was not found in cache", requiredValue)
	}

}

func TestCache_GetAll(t *testing.T) {
	requiredValue := 10

	c := initializeFullCache(requiredValue, nil)

	l := len(c.GetAll())

	if l != requiredValue {
		t.Errorf("Required value %d, got %d", requiredValue, l)
	}
}

func TestCache_Remove(t *testing.T) {
	c := initializeFullCache(10, nil)

	valueToRemove := 5

	c.Remove(valueToRemove)

	if _, exist := c.data[valueToRemove]; exist {
		t.Errorf("Value %d was supposed to be removed from the cache, but it was not", valueToRemove)
	}
}

func TestCache_GetBulk(t *testing.T) {
	c := initializeFullCache(10, nil)
	requiredValues := []int{2, 4, 6}

	results := c.GetBulk(requiredValues)

	for _, i := range requiredValues {
		if n, exist := results[i]; !exist {
			t.Errorf("Expected to see %d, got %d", i, n)
		}
	}
}

func TestCache_Reset(t *testing.T) {
	c := initializeFullCache(10, nil)

	c.Reset()

	l := len(c.data)

	if l != 0 {
		t.Errorf("Expected to have cache of size 0, got %d", l)
	}
}

func TestCache_ForEach(t *testing.T) {
	c := initializeFullCache(10, nil)

	desiredValue := 45
	i := 0

	c.ForEach(func(k, v int) {
		i += v
	})

	if i != desiredValue {
		t.Errorf("Desired value is %d, got %d", desiredValue, i)
	}
}

func TestCache_GetAllAndRemove(t *testing.T) {
	c := initializeFullCache(10, nil)

	d := c.GetAllAndRemove()

	cLen := len(c.data)
	dLen := len(d)

	if dLen != 10 || cLen != 0 {
		t.Errorf("Expected to have 0 elements in cache after GetAllAndRemove() was called and 10 elements returned from it, but received %d elements in cache and %d received from GetAllAndRemove()", cLen, dLen)
	}
}

func TestCache_GetAndRemove(t *testing.T) {
	c := initializeFullCache(10, nil)

	elementToRemove := 5

	c.GetAndRemove(elementToRemove)

	cLen := len(c.data)
	_, exist := c.data[elementToRemove]

	if cLen != 9 || exist {
		t.Errorf("Expected cache length is 9 and presence of the removed element in the cache to be false, got cach length %d and presence %t", cLen, exist)
	}

}

func TestCache_GetRandomSamples(t *testing.T) {
	c := initializeFullCache(10, nil)

	numberOfSamples := 4
	samples := c.GetRandomSamples(numberOfSamples)
	lenSamples := len(samples)

	if lenSamples != numberOfSamples {
		t.Errorf("Expected to have %d samples, got %d", numberOfSamples, lenSamples)
	}

	for k := range samples {
		if _, exist := c.data[k]; !exist {
			t.Errorf("Key %d received from GetRandomSamples() method but it doesn't actually exist in the cache!", k)
		}
	}
}

func TestCache_RemoveBulk(t *testing.T) {
	c := initializeFullCache(10, nil)

	c.RemoveBulk([]int{0, 2, 4, 6, 8})

	expectedLength := 5
	cLen := len(c.data)

	if cLen != expectedLength {
		t.Errorf("Expected cache size is %d, got %d", expectedLength, cLen)
	}
}

func TestNew(t *testing.T) {
	c1 := New[int, int](nil)
	c2 := New[int, int](&Requirements{DefaultTimeout: time.Second * 30})

	c1Len := len(c1.data)
	c2Len := len(c2.data)

	if c1Len > 0 || c2Len > 0 {
		t.Errorf("Expected to have cache sizes of 0 0 0, got %d %d", c1Len, c2Len)
	}

	req1 := c1.Requirements()

	if req1.timeoutInUse {
		t.Errorf("Expected cache1 timeoutInUse to be false, got %t", req1.timeoutInUse)
	}

	req2 := c2.Requirements()

	if !req2.timeoutInUse {
		t.Errorf("Expected cache2 timeoutInUse to be true, got %t", req2.timeoutInUse)
	}

	tm := req2.DefaultTimeout.String()

	if tm != "30s" {
		t.Errorf("Cache2 expected to have DefaultTimeout of 30s, got %s", tm)
	}
}

func TestCopy(t *testing.T) {
	c1 := initializeFullCache(50, &Requirements{DefaultTimeout: time.Second * 30})
	c2 := Copy(c1)

	c2Len := len(c2.data)
	tm := c2.Requirements().DefaultTimeout.String()
	timeoutInUse := c2.Requirements().timeoutInUse

	if c2Len != 50 {
		t.Errorf("Expected cache2 length is 50, got %d", c2Len)
	}

	if tm != "30s" || !timeoutInUse {
		t.Errorf("Expected cache2 to have DefaultTimeout of 30s and timeoutInUse to be true, got %s, %t", tm, timeoutInUse)
	}
}

func TestMerge(t *testing.T) {
	main := initializeFullCache(10, nil)
	secondary := initializeFullCache(20, nil)

	Merge[int, int](main, secondary)

	mainLen := len(main.data)

	if mainLen != 20 {
		t.Errorf("Expected the main cache to have 20 elements in it, got %d", mainLen)
	}
}

func TestMergeAndReset(t *testing.T) {
	main := initializeFullCache(10, nil)
	secondary := initializeFullCache(20, nil)

	MergeAndReset[int, int](main, &secondary)

	mainLen := len(main.data)
	secondaryLen := len(secondary.data)

	if mainLen != 20 {
		t.Errorf("Expected the main cache to have 20 elements in it, got %d", mainLen)
	}

	if secondaryLen != 0 {
		t.Errorf("Expected secondary cache to have 0 items in it, got %d", secondaryLen)
	}
}

func TestCache_Requirements(t *testing.T) {
	c := initializeFullCache(10, &Requirements{DefaultTimeout: time.Millisecond * 500})

	timeoutUsed := c.Requirements().timeoutInUse

	if !timeoutUsed {
		t.Errorf("timeoutInUse expected to be true, got %t", timeoutUsed)
	}

	cLen := c.Count()

	if cLen != 10 {
		t.Errorf("Expected to have 10 items in the cache, got %d", cLen)
	}

	time.Sleep(time.Millisecond * 750)

	cLen = c.Count()

	if cLen != 0 {
		t.Errorf("Expected to have 0 items in the cache, got %d", cLen)
	}
}

func TestEntry_Value(t *testing.T) {
	c := initializeFullCache(0, nil)

	v1 := c.Add(1, 1).Value()
	v2 := c.Add(2, 2).Value()
	v3 := c.Add(3, 3).Value()

	if v1 != 1 || v2 != 2 || v3 != 3 {
		t.Errorf("Expected to have values 1, 2, 3. Got %d, %d, %d", v1, v2, v3)
	}
}

func TestEntry_TimerExist(t *testing.T) {
	c1 := initializeFullCache(1, &Requirements{DefaultTimeout: time.Second * 30})
	c2 := initializeFullCache(1, nil)

	c1Exist := c1.Add(2, 2).TimerExist()
	c2Exist := c2.Add(2, 2).TimerExist()

	if !c1Exist || c2Exist {
		t.Errorf("Expected TimerExist method to return true from cache1 and false from cache2, got %t, %t", c1Exist, c2Exist)
	}
}

func TestCache_GetEntry(t *testing.T) {
	c := initializeFullCache(10, nil)

	val := c.GetEntry(2).Value()

	if val != 2 {
		t.Errorf("Expected to get value %d. Got %d", 2, val)
	}
}

func TestCache_GetAndRemoveEntry(t *testing.T) {
	c := initializeFullCache(10, nil)

	val := c.GetAndRemoveEntry(2).Value()

	if val != 2 {
		t.Errorf("Expected to get value %d. Got %d", 2, val)
	}

	if c.Exist(2) {
		t.Errorf("Key %d in cache shouldn't exist, but it does!", 2)
	}
}

func TestCache_AddWithTimeout(t *testing.T) {
	c := initializeFullCache(0, nil)

	val := 1

	c.AddWithTimeout(val, val, time.Millisecond*500)

	if !c.Exist(val) {
		t.Errorf("Value with key %d should exist in the cache, but it does not!", val)
	}

	time.Sleep(time.Millisecond * 1000)

	if c.Exist(val) {
		t.Errorf("Value with key %d should NOT exist in the cache, but it does!", val)
	}
}

func TestCache_AddTimer(t *testing.T) {
	c := initializeFullCache(10, nil)

	e1 := c.GetEntry(2)

	if e1.TimerExist() {
		t.Errorf("Timer does not suppose to exist, but it does!")
	}

	c.AddTimer(2, time.Second*30)

	e2 := c.GetEntry(2)

	if !e2.TimerExist() {
		t.Errorf("Timer suppose to exist, but it does not!")
	}
}

func TestEntry_StopTimer(t *testing.T) {
	c := initializeFullCache(10, &Requirements{DefaultTimeout: time.Millisecond * 250})

	e := c.GetEntry(1)

	if e == nil {
		t.Errorf("Expected to have entry using key %d, got <nil>", 1)
		return
	}

	e.StopTimer()

	time.Sleep(time.Millisecond * 500)

	if !c.Exist(1) {
		t.Errorf("Entry with key %d should be preset, but it is not!", 1)
	}
}

func TestEntry_ResetTimer(t *testing.T) {
	c := initializeFullCache(10, &Requirements{DefaultTimeout: time.Millisecond * 250})

	if !c.Exist(1) || !c.Exist(2) {
		t.Errorf("Both, entry 1 and 2 should be present in the cache, but one or both are not!")
	}

	c.GetEntry(1).ResetTimer(time.Millisecond * 500)

	time.Sleep(time.Millisecond * 350)

	if c.Exist(2) {
		t.Errorf("Entry with key 2 should not exist in cache, but it does!")
	}

	if !c.Exist(1) {
		t.Errorf("Entry with key 1 should exist in the cache, but it does not!")
	}
}

//===========[BENCHMARKS]====================================================================================================

func BenchmarkCache_AddTimer(b *testing.B) {
	c := initializeFullCache(10, nil)

	t := time.Second * 90

	for n := 0; n < b.N; n++ {
		c.AddTimer(2, t)
	}
}

func BenchmarkAdd(b *testing.B) {
	c := initializeFullCache(0, nil)

	for n := 0; n < b.N; n++ {
		c.Add(n, n)
	}
}

func BenchmarkAddBulk(b *testing.B) {
	c := initializeFullCache(0, nil)

	for n := 0; n < b.N; n++ {
		c.AddBulk(map[int]int{
			n: n,
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	c := initializeFullCache(0, nil)

	for n := 0; n < b.N; n++ {
		c.Remove(n)
	}
}

func BenchmarkRemoveBulk(b *testing.B) {
	c := initializeFullCache(0, nil)

	for n := 0; n < b.N; n++ {
		c.RemoveBulk([]int{n, n + 1, n + 2})
	}
}

func BenchmarkExist(b *testing.B) {
	c := initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		c.Exist(1)
	}
}

func BenchmarkGet(b *testing.B) {
	c := initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		c.Get(1)
	}
}

func BenchmarkGetBulk(b *testing.B) {
	c := initializeFullCache(1, nil)

	for n := 0; n < b.N; n++ {
		c.GetBulk([]int{0})
	}
}

func BenchmarkGetAndRemove(b *testing.B) {
	c := initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		c.GetAndRemove(1)
	}
}

func BenchmarkGetAll(b *testing.B) {
	c := initializeFullCache(1, nil)

	for n := 0; n < b.N; n++ {
		c.GetAll()
	}
}

func BenchmarkCount(b *testing.B) {
	c := initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		c.Count()
	}
}

func BenchmarkReset(b *testing.B) {
	var c = initializeFullCache(10, nil)

	for n := 0; n < b.N; n++ {
		c.Reset()
	}
}

func BenchmarkForEach(b *testing.B) {
	cache := initializeFullCache(1, nil)

	for n := 0; n < b.N; n++ {
		cache.ForEach(func(key, val int) {})
	}
}

func BenchmarkCopy(b *testing.B) {
	var c1 = initializeFullCache(1, nil)

	for n := 0; n < b.N; n++ {
		Copy[int, int](c1)
	}

}

func BenchmarkMerge(b *testing.B) {
	var c1 = initializeFullCache(1, nil)
	var c2 = initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		Merge[int, int](c1, c2)
	}

}

func BenchmarkMergeAndReset(b *testing.B) {
	var c1 = initializeFullCache(1, nil)
	var c2 = initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		MergeAndReset[int, int](c1, c2)
	}

}
