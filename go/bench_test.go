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

func makeThreadSafeLRU(n int) LRUCache {
	a := NewLRUCache(n)
	nums := nrand(n)
	for _, v := range nums {
		a.Add(v, v)
	}
	return a
}

func makeThreadUnsafeLRU(n int) LRUCache {
	a := NewThreadUnsafeLRUCache(n)
	nums := nrand(n)
	for _, v := range nums {
		a.Add(v, v)
	}
	return a
}

func benchAdd(b *testing.B, lru LRUCache) {
	nums := nrand(b.N)
	b.ResetTimer()
	for _, v := range nums {
		lru.Add(v, v)
	}
}

func benchSize(b *testing.B, lru LRUCache) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Size()
	}
}

func benchFind(b *testing.B, lru LRUCache) {
	nums := nrand(b.N)
	b.ResetTimer()
	for _, v := range nums {
		lru.Find(v)
	}
}

func benchIterator(b *testing.B, lru LRUCache) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Iterator(true)
	}
}

func benchIter(b *testing.B, lru LRUCache) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Iter(true)
	}
}

// benchmark thread safe lru

func BenchmarkThreadSafeLRU_Add(b *testing.B) {
	a := NewLRUCache(100)
	benchAdd(b, a)
}

func BenchmarkThreadSafeLRU_Size(b *testing.B) {
	a := NewLRUCache(100)
	benchSize(b, a)
}

func BenchmarkThreadSafeLRU_Find(b *testing.B) {
	a := makeThreadSafeLRU(b.N / 2)
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Find2(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Iterator(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchIterator(b, a)
}

func BenchmarkThreadSafeLRU_Itera(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchIter(b, a)
}

// benchmark thread unsafe lru

func BenchmarkThreadUnsafeLRU_Add(b *testing.B) {
	a := NewThreadUnsafeLRUCache(100)
	benchAdd(b, a)
}

func BenchmarkThreadUnsafeLRU_Size(b *testing.B) {
	a := NewThreadUnsafeLRUCache(100)
	benchSize(b, a)
}

func BenchmarkThreadUnsafeLRU_Find(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N / 2)
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Find2(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Iterator(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchIterator(b, a)
}

func BenchmarkThreadUnsafeLRU_Iter(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchIter(b, a)
}
