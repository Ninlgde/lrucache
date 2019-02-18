package lru

import (
	"github.com/deckarep/golang-set"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

const N = 1000
const CAP = 1001

func TestThreadSafeLRU_Add_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
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

	Assert(a.Size() == CAP/2, t)
	findCount := 0
	for _, i := range ints {
		if a.Find(i) != nil {
			findCount++
		}
	}

	Assert(findCount == CAP/2, t)
}

func TestThreadSafeLRU_Find_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
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

	Assert(a.Size() == CAP/2, t)

	findCount := 0
	for _, i := range ints {
		if result.Contains(i) {
			findCount++
		}
	}

	Assert(findCount == CAP/2, t)
}

func TestThreadSafeLRU_Remove_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
	ints := rand.Perm(N)

	for i := 0; i < len(ints); i++ {
		a.Add(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	result := mapset.NewSet() // concurrent set
	for _, i := range ints {
		go func(i int) {
			result.Add(a.Remove(i))
			wg.Done()
		}(i)
	}
	wg.Wait()

	Assert(a.Size() == 0, t)

	removeCount := 0
	for _, i := range ints {
		if result.Contains(i) {
			removeCount++
		}
		Assert(a.Find(i) == nil, t)
	}

	Assert(removeCount == CAP/2, t)
}

func TestThreadSafeLRU_Iterator_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
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

	// 理论上 所有iterator里的元素相加应该是 大于 (1+..+1000)/2 小于 500*1000的
	count := 0
	for _, iter := range resultIterator {
		count += len(iter.C)
	}

	Assert(count < N*N/2, t)
	Assert(count > (1+N)*N/4, t)
}

func TestThreadSafeLRU_Iter_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
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

	// 理论上 所有iterator里的元素相加应该是 大于 (1+..+1000)/2 小于 500*1000的
	count := 0
	for _, iter := range resultIters {
		count += len(iter)
	}

	Assert(count < N*N/2, t)
	Assert(count > (1+N)*N/4, t)
}

func TestThreadSafeLRU(t *testing.T) {
	runtime.GOMAXPROCS(2)

	a := NewLRUCache(CAP / 2)
	rand.Seed(time.Now().UnixNano()) // 尽量随机
	ints := rand.Perm(N)

	for i := 0; i < len(ints); i++ {
		a.Add(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	//result := mapset.NewSet() // concurrent set
	addCount := 0
	findResult := make([]lruValue, 0, N)
	removeResult := make([]lruValue, 0, N)
	iterResult := make([]<-chan lruPair, 0, N)
	for _, i := range ints {
		go func(i int) {
			switch rand.Int() % 4 {
			case 0:
				a.Add(i, i)
				addCount++
			case 1:
				findResult = append(findResult, a.Find(i))
			case 2:
				removeResult = append(removeResult, a.Remove(i))
			case 3:
				iterResult = append(iterResult, a.Iter(true))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// assert add
	Assert(addCount > 200 && addCount < 300, t)

	// assert find
	findCount := 0
	for _, f := range findResult {
		if f != nil {
			findCount++
		}
	}
	Assert(findCount > 0 && findCount < len(findResult), t)

	// assert remove
	removeCount := 0
	for _, f := range removeResult {
		if f != nil {
			removeCount++
		}
	}
	Assert(removeCount > 0 && removeCount < len(removeResult), t)

	// assert iter
	iterCount := 0
	for _, iter := range iterResult {
		iterCount += len(iter)
	}
	Assert(iterCount < len(iterResult)*N/2, t)      // 小于 500*250
	Assert(iterCount > len(iterResult)*(N/2-50), t) // 概率上 不会离500太远

	Assert(a.Size() > 450, t)
}
