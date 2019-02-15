package lru

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const N = 1000
const CAP = 1001

func TestThreadSafeLRU_Add_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP)
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			a.Add(i, i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	Assert(a.Size() == N, t)
	for _, i := range ints {
		Assert(a.Find(i) == i, t)
	}

	Assert(a.Find(N) == nil, t)
}

func TestThreadSafeLRU_Find_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP)
	ints := rand.Perm(N)

	for i := 0; i < len(ints); i++ {
		a.Add(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	result := mapset.NewSet() // concurrent set
	for _, i := range ints {
		go func(i int) {
			result.Add(a.Find(i))
			wg.Done()
		}(i)
	}
	wg.Wait()

	Assert(a.Size() == N, t)
	Assert(result.Cardinality() == N, t)

	for _, i := range ints {
		Assert(result.Contains(i), t)
	}

	Assert(!result.Contains(N), t)
}

func TestThreadSafeLRU_Iterator_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP)
	ints := rand.Perm(N)

	resultIterator := make([]*Iterator, 0, 1000)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			a.Add(i, i)
			go func() {
				iter := a.Iterator(true)
				resultIterator = append(resultIterator, iter)
			}()
			wg.Done()
		}(i)
	}
	wg.Wait()

	// 理论上 所有iterator里的元素相加应该是 大于 (1+..+1000)/2 小于 1000*1000的
	count := 0
	for _, iter := range resultIterator {
		count += len(iter.C)
	}

	fmt.Println(count)
	Assert(count < N*N, t)
	Assert(count > (1+N)*N/4, t)
}

func TestThreadSafeLRU_Iter_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP)
	ints := rand.Perm(N)

	resultIters := make([]<-chan lruPair, 0, 1000)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			a.Add(i, i)
			go func() {
				iter := a.Iter(true)
				resultIters = append(resultIters, iter)
			}()
			wg.Done()
		}(i)
	}
	wg.Wait()

	// 理论上 所有iterator里的元素相加应该是 大于 (1+..+1000)/2 小于 1000*1000的
	count := 0
	for _, iter := range resultIters {
		count += len(iter)
	}

	fmt.Println(count)
	Assert(count < N*N, t)
	Assert(count > (1+N)*N/4, t)
}

func TestSet(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := mapset.NewSet()
	ints := rand.Perm(N)

	//cs := make([]<-chan interface{}, 0)
	var wg sync.WaitGroup
	wg.Add(len(ints) * 2)
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
		go func(i int) {
			//iter := s.Iter()
			//cs = append(cs, iter)
			s.Remove(i)
			wg.Done()
		}(i)
	}
	//for i := 0; i < len(ints); i++ {
	//}
	wg.Wait()
	fmt.Println(s.Cardinality())
}
