package cacheMachine

import (
	"testing"
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

	for k, _ := range samples {
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