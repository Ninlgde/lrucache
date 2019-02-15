package lru

import (
	"math/rand"
	"testing"
)

func nrand(n int) []int {
	i := make([]int, n)
	for ind := range i {
		i[ind] = rand.Int()
	}
	return i
}

func benchAdd(b *testing.B, lru LRUCache) {
	nums := nrand(b.N)
	b.ResetTimer()
	for _, v := range nums {
		lru.Add(v, v)
	}
}

func BenchmarkThreadSafeLRU_Add(b *testing.B) {
	a := NewLRUCache(100)
	benchAdd(b, a)
}

func BenchmarkThreadUnsafeLRU_Add(b *testing.B) {
	a := NewThreadUnsafeLRUCache(100)
	benchAdd(b, a)
}

func benchSize(b *testing.B, lru LRUCache) {
	for i := 0; i < b.N; i++ {
		lru.Size()
	}
}

func BenchmarkThreadSafeLRU_Size(b *testing.B) {
	a := NewLRUCache(100)
	benchSize(b, a)
}

func BenchmarkThreadUnsafeLRU_Size(b *testing.B) {
	a := NewThreadUnsafeLRUCache(100)
	benchSize(b, a)
}

func benchFind(b *testing.B, lru LRUCache) {
	nums := nrand(b.N)
	b.ResetTimer()
	for _, v := range nums {
		lru.Add(v, v)
	}
}

func BenchmarkThreadSafeLRU_Find(b *testing.B) {
	a := NewLRUCache(b.N / 2)
	nums := nrand(b.N / 2)
	for _, v := range nums {
		a.Add(v, v)
	}
	b.ResetTimer()
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Find2(b *testing.B) {
	a := NewLRUCache(b.N)
	nums := nrand(b.N)
	for _, v := range nums {
		a.Add(v, v)
	}
	b.ResetTimer()
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Find(b *testing.B) {
	a := NewThreadUnsafeLRUCache(b.N / 2)
	nums := nrand(b.N / 2)
	for _, v := range nums {
		a.Add(v, v)
	}
	b.ResetTimer()
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Find2(b *testing.B) {
	a := NewThreadUnsafeLRUCache(b.N)
	nums := nrand(b.N)
	for _, v := range nums {
		a.Add(v, v)
	}
	b.ResetTimer()
	benchFind(b, a)
}
