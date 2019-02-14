package lru

import "testing"

func Assert(con bool, t *testing.T) {
	if !con {
		t.Error("Assert error")
	}
}

func AssertPairList(a, b []lruPair, t *testing.T) {
	if len(a) != len(b) {
		t.Error("AssertPairList len error")
		return
	}
	for i, v := range a {
		if v != b[i] {
			t.Errorf("%v != %v\n", v, b[i])
		}
	}
}

// thread safe lru interface tests
// multi thread test see threadsafe_test.go
func TestNewLRUCache(t *testing.T) {
	a := NewLRUCache(10)

	Assert(a.Size() == 0, t)
}

func TestThreadSafeLRU_Add(t *testing.T) {
	a := NewLRUCache(10)

	v3 := &struct{}{}
	a.Add(1, 1)
	a.Add(2, "two")
	a.Add(3, v3)

	Assert(a.Size() == 3, t)
}

func TestThreadSafeLRU_Add_Core(t *testing.T) {
	a := NewLRUCache(10)

	v3 := &struct{}{}
	a.Add(1, 1)
	a.Add(2, "two")
	a.Add(3, v3)

	Assert(a.Size() == 3, t)

	except := []lruPair{lruPair{3, v3}, lruPair{2, "two"}, lruPair{1, 1}}
	result := make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		result = append(result, p)
	}
	AssertPairList(except, result, t)

	a.Add(2, "2")
	except = []lruPair{lruPair{2, "2"}, lruPair{3, v3}, lruPair{1, 1}}
	result = make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		result = append(result, p)
	}
	AssertPairList(except, result, t)
}

func TestThreadSafeLRU_Find(t *testing.T) {
	a := NewLRUCache(10)

	v3 := &struct{}{}
	a.Add(1, 1)
	a.Add(2, "two")
	a.Add(3, v3)

	Assert(a.Find(1) == 1, t)
	Assert(a.Find(2) == "two", t)
	Assert(a.Find(3) == v3, t)
	Assert(a.Find(4) == nil, t)
	// find can not inc size
	Assert(a.Size() == 3, t)
}

func TestThreadSafeLRU_Find_Core(t *testing.T) {
	a := NewLRUCache(10)

	v3 := &struct{}{}
	a.Add(1, 1)
	a.Add(2, "two")
	a.Add(3, v3)

	Assert(a.Size() == 3, t)

	a.Find(2)
	except := []lruPair{lruPair{2, "two"}, lruPair{3, v3}, lruPair{1, 1}}
	result := make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		result = append(result, p)
	}
	AssertPairList(except, result, t)

	Assert(a.Find(4) == nil, t)
	// not found can not change order
	a.Find(2)
	except = []lruPair{lruPair{2, "two"}, lruPair{3, v3}, lruPair{1, 1}}
	result = make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		result = append(result, p)
	}
	AssertPairList(except, result, t)
}

func TestThreadSafeLRU_Create(t *testing.T) {
	a := NewLRUCache(10)
	for i := 0; i <= 10; i++ {
		a.Add(i, i) // add 11 elem
	}

	// 10 >> 2 = 2; 10 - 2 + 1 = 9
	Assert(a.Size() == 9, t)

	b := NewLRUCache(100)
	for i := 0; i <= 100; i++ {
		b.Add(i, i) // add 101 elem
	}

	// 100 >> 2 = 25; 100 - 25 + 1 = 76
	Assert(b.Size() == 76, t)
}

func TestThreadSafeLRU_Iter(t *testing.T) {
	a := NewLRUCache(10)
	exceptt := make([]lruPair, 10)
	exceptf := make([]lruPair, 10)
	for i := 0; i < 10; i++ {
		a.Add(i, i) // add 10 elem
		exceptt[i] = lruPair{i, i}
		exceptf[9-i] = lruPair{i, i}
	}

	resultt := make([]lruPair, 0, 10)
	for p := range a.Iter(true) {
		resultt = append(resultt, p)
	}

	resultf := make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		resultf = append(resultf, p)
	}

	AssertPairList(exceptt, resultt, t)
	AssertPairList(exceptf, resultf, t)
}

func TestThreadSafeLRU_Iterator(t *testing.T) {
	a := NewLRUCache(10)
	exceptt := make([]lruPair, 10)
	exceptf := make([]lruPair, 10)
	for i := 0; i < 10; i++ {
		a.Add(i, i) // add 10 elem
		exceptt[i] = lruPair{i, i}
		exceptf[9-i] = lruPair{i, i}
	}

	iterator := a.Iterator(true)
	resultt := make([]lruPair, 0, 10)
	for p := range iterator.C {
		resultt = append(resultt, p)
		if p.k.(int) == 5 {
			iterator.Stop()
		}
	}

	iterator = a.Iterator(false)
	resultf := make([]lruPair, 0, 10)
	for p := range iterator.C {
		resultf = append(resultf, p)
		if p.k.(int) == 5 {
			iterator.Stop()
		}
	}

	Assert(len(resultt) == 6, t)
	Assert(len(resultf) == 5, t)

	AssertPairList(exceptt[:6], resultt, t)
	AssertPairList(exceptf[:5], resultf, t)
}

// thread unsafe lru interface tests
func TestNewThreadUnsafeLRUCache(t *testing.T) {
	a := NewThreadUnsafeLRUCache(10)

	Assert(a.Size() == 0, t)
}

// same as thread safe lru
//func TestThreadUnsafeLRU_Add(t *testing.T)
//
//func TestThreadUnsafeLRU_Add_Core(t *testing.T)
//
//func TestThreadUnsafeLRU_Find(t *testing.T)
//
//func TestThreadUnsafeLRU_Find_Core(t *testing.T)
//
//func TestThreadUnsafeLRU_Create(t *testing.T)

func TestThreadUnsafeLRU_Iter(t *testing.T) {
	a := NewThreadUnsafeLRUCache(10)
	exceptt := make([]lruPair, 10)
	exceptf := make([]lruPair, 10)
	for i := 0; i < 10; i++ {
		a.Add(i, i) // add 10 elem
		exceptt[i] = lruPair{i, i}
		exceptf[9-i] = lruPair{i, i}
	}

	resultt := make([]lruPair, 0, 10)
	for p := range a.Iter(true) {
		resultt = append(resultt, p)
	}

	resultf := make([]lruPair, 0, 10)
	for p := range a.Iter(false) {
		resultf = append(resultf, p)
	}

	AssertPairList(exceptt, resultt, t)
	AssertPairList(exceptf, resultf, t)
}

func TestThreadUnsafeLRU_Iterator(t *testing.T) {
	a := NewThreadUnsafeLRUCache(10)
	exceptt := make([]lruPair, 10)
	exceptf := make([]lruPair, 10)
	for i := 0; i < 10; i++ {
		a.Add(i, i) // add 10 elem
		exceptt[i] = lruPair{i, i}
		exceptf[9-i] = lruPair{i, i}
	}

	iterator := a.Iterator(true)
	resultt := make([]lruPair, 0, 10)
	for p := range iterator.C {
		resultt = append(resultt, p)
		if p.k.(int) == 5 {
			iterator.Stop()
		}
	}

	iterator = a.Iterator(false)
	resultf := make([]lruPair, 0, 10)
	for p := range iterator.C {
		resultf = append(resultf, p)
		if p.k.(int) == 5 {
			iterator.Stop()
		}
	}

	Assert(len(resultt) == 6, t)
	Assert(len(resultf) == 5, t)

	AssertPairList(exceptt[:6], resultt, t)
	AssertPairList(exceptf[:5], resultf, t)
}

func TestThreadUnsafeLRU_Iterator2(t *testing.T) {
	a := NewThreadUnsafeLRUCache(10)
	iter := a.Iterator(true)
	c := make(chan lruPair)
	for p := range iter.C {
		c <- p
	}
}
