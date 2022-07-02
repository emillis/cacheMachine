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
