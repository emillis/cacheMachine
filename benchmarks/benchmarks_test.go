package benchmarks

import (
	"github.com/emillis/cacheMachine"
	"testing"
)

//===========[FUNCTIONALITY]====================================================================================================

func populateCache(n int, c cacheMachine.Cache[int, int]) {
	for i := 0; i < n; i++ {
		c.Add(i, i)
	}
}

func initializeFullCache(n int, r *cacheMachine.Requirements) cacheMachine.Cache[int, int] {
	c := cacheMachine.New[int, int](r)

	for i := 0; i < n; i++ {
		c.Add(i, i)
	}

	return c
}

//===========[BENCHMARKS]====================================================================================================

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
	c := cacheMachine.New[int, int](nil)
	populateCache(5, c)

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
	var c = cacheMachine.New[int, int](nil)

	for n := 0; n < b.N; n++ {
		c.Reset()
	}
}

func BenchmarkForEach(b *testing.B) {
	cache := cacheMachine.New[int, int](nil)

	populateCache(1, cache)

	for n := 0; n < b.N; n++ {
		cache.ForEach(func(key, val int) {})
	}
}

func BenchmarkCopy(b *testing.B) {
	var c1 = initializeFullCache(1, nil)

	for n := 0; n < b.N; n++ {
		cacheMachine.Copy[int, int](c1)
	}

}

func BenchmarkMerge(b *testing.B) {
	var c1 = initializeFullCache(1, nil)
	var c2 = initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		cacheMachine.Merge[int, int](c1, c2)
	}

}

func BenchmarkMergeAndReset(b *testing.B) {
	var c1 = initializeFullCache(1, nil)
	var c2 = initializeFullCache(2, nil)

	for n := 0; n < b.N; n++ {
		cacheMachine.MergeAndReset[int, int](c1, c2)
	}

}
