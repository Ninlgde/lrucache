package lru

import (
	"sync"
)

/**
LRU 缓存
由一个map和一个双向链表组成
可将查找、添加等操作的时间复杂度较少到O(1) (理论上，取决于map的实现)
*/
type threadSafeLRU struct {
	c            *threadUnsafeLRU
	sync.RWMutex // 协程锁
}

func newThreadSafeLRU() *threadSafeLRU {
	return &threadSafeLRU{}
}

/**
创建缓存
cap: 容量，缓存最多存多少数据
*/
func (cache *threadSafeLRU) Create(cap int) {
	cache.c = newThreadUnsafeLRU()
	cache.c.Create(cap)
}

/**
添加一个元素
k: key
v: value

cost: O(1)(base on map's implement)
*/
func (cache *threadSafeLRU) Add(k lruKey, v lruValue) {
	cache.Lock()
	defer cache.Unlock()
	cache.c.Add(k, v)
}

/**
查找一个元素
k: key
return: value or nil

cost: O(1)(base on map's implement)
*/
func (cache *threadSafeLRU) Find(k lruKey) lruValue {
	cache.Lock()
	defer cache.Unlock()
	return cache.c.Find(k)
}

/**
当前缓存大小
return: size of cache

cost: O(1)
*/
func (cache *threadSafeLRU) Size() int {
	cache.RLock()
	defer cache.RUnlock()
	return cache.c.Size()
}

/**
删除一个元素
k: key
return: value or nil

cost: O(1)
*/
func (cache *threadSafeLRU) Remove(k lruKey) lruValue {
	cache.Lock()
	defer cache.Unlock()
	return cache.c.Remove(k)
}

/**
遍历缓存中所有的数据的迭代器
reverse: 是否翻转 true = 正序 false = 倒序(默认，淘汰的是从头部，所以从后往前是默认)
return: 迭代器 func
*/
func (cache *threadSafeLRU) Iterator(reverse bool) *Iterator {
	iterator, ch, stopCh := newIterator(cache.c.cap)
	go func() {
		cache.RLock()
		defer cache.RUnlock()
		if cache.c.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.c.head.next
		LT:
			for p != cache.c.tail && p != nil {
				select {
				case <-stopCh:
					break LT
				case ch <- lruPair{p.key, p.value}:
				}
				p = p.next
			}
		} else {
			p := cache.c.tail.prev
		LF:
			for p != cache.c.head {
				select {
				case <-stopCh:
					break LF
				case ch <- lruPair{p.key, p.value}:
				}
				p = p.prev
			}
		}
		close(ch)
	}()
	return iterator
}

func (cache *threadSafeLRU) Iter(reverse bool) <-chan lruPair {
	ch := make(chan lruPair, cache.c.cap) // 这里需要设置channel大小
	go func() {
		cache.RLock()
		defer cache.RUnlock()
		if cache.c.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.c.head.next
			for p != cache.c.tail {
				ch <- lruPair{p.key, p.value}
				p = p.next
			}
		} else {
			p := cache.c.tail.prev
			for p != cache.c.head {
				ch <- lruPair{p.key, p.value}
				p = p.prev
			}
		}
		close(ch)
	}()

	return ch
}
