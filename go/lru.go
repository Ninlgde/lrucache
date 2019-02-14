package lru

// 值的类型 类似void*
type lruValue interface{}

// key的类型  类似void*
type lruKey interface{}

type lruPair struct {
	k lruKey
	v lruValue
}

type LRUCache interface {
	Create(cap int)
	Add(k lruKey, v lruValue)
	Size() int
	Find(k lruKey) lruValue
	Iterator(reverse bool) *Iterator
	Iter(reverse bool) <-chan lruPair
}

func NewLRUCache(cap int) LRUCache {
	lru := newThreadSafeLRU()
	lru.Create(cap)
	return lru
}

func NewThreadUnsafeLRUCache(cap int) LRUCache {
	lru := newThreadUnsafeLRU()
	lru.Create(cap)
	return lru
}
