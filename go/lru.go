package lru

// 值的类型 类似void*(clang)
type lruValue interface{}

// key的类型  类似void*(clang)
type lruKey interface{}

type lruPair struct {
	k lruKey
	v lruValue
}

// LRU Cache
// a fast lru cache implement by ninlgde
type LRUCache interface {
	// create a lru cache with cap(capacity)
	Create(cap int)

	// add key and value to lru cache
	Add(k lruKey, v lruValue)

	// get the size of lru cache
	Size() int

	// find key in lru cache
	// if find, move the node to the tail
	Find(k lruKey) lruValue

	// remove a key in lru cache
	Remove(k lruKey) lruValue

	// the iterators
	Iterator(reverse bool) *Iterator
	Iter(reverse bool) <-chan lruPair
}

// new a thread safe lru cache
func NewLRUCache(cap int) LRUCache {
	lru := newThreadSafeLRU()
	lru.Create(cap)
	return lru
}

// new a thread unsafe lru cache
func NewThreadUnsafeLRUCache(cap int) LRUCache {
	lru := newThreadUnsafeLRU()
	lru.Create(cap)
	return lru
}
