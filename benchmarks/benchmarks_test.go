package benchmarks

import (
	"github.com/emillis/cacheMachine"
	"testing"
)

var cache = cacheMachine.New[int, int](nil)

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
	cache.Add(7, 8)

	for n := 0; n < b.N; n++ {
		cache.GetBulk([]int{0})
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

func BenchmarkGetRandomSamples(b *testing.B) {
	var c = cacheMachine.New[int, int](map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
	})

	for n := 0; n < b.N; n++ {
		c.GetRandomSamples(1)
	}

}
