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
	nums := nrand(b.N * 2)
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

func benchRemove(b *testing.B, lru LRUCache) {
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
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := NewLRUCache(n)
	benchAdd(b, a)
}

func BenchmarkThreadSafeLRU_Add2(b *testing.B) {
	a := NewLRUCache(b.N)
	benchAdd(b, a)
}

func BenchmarkThreadSafeLRU_Add3(b *testing.B) {
	a := NewLRUCache(100)
	benchAdd(b, a)
}

func BenchmarkThreadSafeLRU_Size(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := NewLRUCache(n)
	benchSize(b, a)
}

func BenchmarkThreadSafeLRU_Size2(b *testing.B) {
	a := NewLRUCache(b.N)
	benchSize(b, a)
}

func BenchmarkThreadSafeLRU_Find(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := makeThreadSafeLRU(n)
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Find2(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Find3(b *testing.B) {
	a := makeThreadSafeLRU(100)
	benchFind(b, a)
}

func BenchmarkThreadSafeLRU_Remove(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := makeThreadSafeLRU(n)
	benchRemove(b, a)
}

func BenchmarkThreadSafeLRU_Remove2(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchRemove(b, a)
}

func BenchmarkThreadSafeLRU_Remove3(b *testing.B) {
	a := makeThreadSafeLRU(100)
	benchRemove(b, a)
}

func BenchmarkThreadSafeLRU_Iterator(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchIterator(b, a)
}

func BenchmarkThreadSafeLRU_Iterator2(b *testing.B) {
	a := makeThreadSafeLRU(100)
	benchIterator(b, a)
}

func BenchmarkThreadSafeLRU_Itera(b *testing.B) {
	a := makeThreadSafeLRU(b.N)
	benchIter(b, a)
}

func BenchmarkThreadSafeLRU_Itera2(b *testing.B) {
	a := makeThreadSafeLRU(100)
	benchIter(b, a)
}

// benchmark thread unsafe lru

func BenchmarkThreadUnsafeLRU_Add(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := NewThreadUnsafeLRUCache(n)
	benchAdd(b, a)
}

func BenchmarkThreadUnsafeLRU_Add2(b *testing.B) {
	a := NewThreadUnsafeLRUCache(b.N)
	benchAdd(b, a)
}

func BenchmarkThreadUnsafeLRU_Add3(b *testing.B) {
	a := NewThreadUnsafeLRUCache(100)
	benchAdd(b, a)
}

func BenchmarkThreadUnsafeLRU_Size(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := NewThreadUnsafeLRUCache(n)
	benchSize(b, a)
}

func BenchmarkThreadUnsafeLRU_Size2(b *testing.B) {
	a := NewThreadUnsafeLRUCache(b.N)
	benchSize(b, a)
}

func BenchmarkThreadUnsafeLRU_Find(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := makeThreadUnsafeLRU(n)
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Find2(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Find3(b *testing.B) {
	a := makeThreadUnsafeLRU(100)
	benchFind(b, a)
}

func BenchmarkThreadUnsafeLRU_Remove(b *testing.B) {
	n := b.N / 2
	if n == 0 {
		n = 1
	}
	a := makeThreadUnsafeLRU(n)
	benchRemove(b, a)
}

func BenchmarkThreadUnsafeLRU_Remove2(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchRemove(b, a)
}

func BenchmarkThreadUnsafeLRU_Remove3(b *testing.B) {
	a := makeThreadUnsafeLRU(100)
	benchRemove(b, a)
}

func BenchmarkThreadUnsafeLRU_Iterator(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchIterator(b, a)
}

func BenchmarkThreadUnsafeLRU_Iterator2(b *testing.B) {
	a := makeThreadUnsafeLRU(100)
	benchIterator(b, a)
}

func BenchmarkThreadUnsafeLRU_Iter(b *testing.B) {
	a := makeThreadUnsafeLRU(b.N)
	benchIter(b, a)
}

func BenchmarkThreadUnsafeLRU_Iter2(b *testing.B) {
	a := makeThreadUnsafeLRU(100)
	benchIter(b, a)
}
