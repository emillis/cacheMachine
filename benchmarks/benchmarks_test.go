package benchmarks

import (
	"github.com/emillis/cacheMachine"
	"testing"
)

//===========[STATIC/CACHE]====================================================================================================

var cache = cacheMachine.New[int, int](nil)

//===========[FUNCTIONALITY]====================================================================================================

func populateCache(n int, c cacheMachine.Cache[int, int]) {
	for i := 0; i < n; i++ {
		c.Add(i, i)
	}
}

//===========[BENCHMARKS]====================================================================================================

func BenchmarkAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cache.Add(n, n)
	}
}

func BenchmarkAddBulk(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cache.AddBulk(map[int]int{
			n: n,
			n: n,
			n: n,
			n: n,
			n: n,
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cache.Remove(n)
	}
}

func BenchmarkRemoveBulk(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cache.RemoveBulk([]int{n, n + 1, n + 2})
	}
}

func BenchmarkExist(b *testing.B) {
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.Exist(7)
	}
}

func BenchmarkGet(b *testing.B) {
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.Get(7)
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
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.GetAndRemove(7)
	}
}

func BenchmarkGetAll(b *testing.B) {
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.GetAll()
	}
}

func BenchmarkCount(b *testing.B) {
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.Count()
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
	var c1 = cacheMachine.New[int, int](map[int]int{1: 1})

	for n := 0; n < b.N; n++ {
		cacheMachine.Copy[int, int](c1)
	}

}

func BenchmarkMerge(b *testing.B) {
	var c1 = cacheMachine.New[int, int](map[int]int{1: 1})
	var c2 = cacheMachine.New[int, int](map[int]int{2: 2})

	for n := 0; n < b.N; n++ {
		cacheMachine.Merge[int, int](c1, c2)
	}

}

func BenchmarkMergeAndReset(b *testing.B) {
	var c1 = cacheMachine.New[int, int](map[int]int{1: 1})
	var c2 = cacheMachine.New[int, int](map[int]int{2: 2})

	for n := 0; n < b.N; n++ {
		cacheMachine.MergeAndReset[int, int](c1, c2)
	}

}
