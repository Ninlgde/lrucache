package lru

import "C"
import (
	"sync"
)

/**
LRU 缓存
由一个map和一个双向链表组成
可将查找、添加等操作的时间复杂度较少到O(1) (理论上，取决于map的实现)
*/
type threadSafeLRU struct {
	C            *threadUnsafeLRU
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
	cache.C = newThreadUnsafeLRU()
	cache.C.Create(cap)
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
	cache.C.Add(k, v)
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
	return cache.C.Find(k)
}

/**
当前缓存大小
return: size of cache

cost: O(1)

测试方法，正式不支持
*/
func (cache *threadSafeLRU) Size() int {
	cache.RLock()
	defer cache.RUnlock()
	return cache.C.Size()
}

/**
遍历缓存中所有的数据的迭代器
reverse: 是否翻转 true = 正序 false = 倒序(默认，淘汰的是从头部，所以从后往前是默认)
return: 迭代器 func

测试方法，正式不支持
*/
func (cache threadSafeLRU) Iterator(reverse bool) *Iterator {
	iterator, ch, stopCh := newIterator()
	go func() {
		cache.RLock()
		defer cache.RUnlock()
		if cache.C.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.C.head.next
		LT:
			for p != cache.C.tail && p != nil {
				select {
				case <-stopCh:
					break LT
				case ch <- lruPair{p.key, p.value}:
				}
				p = p.next
			}
		} else {
			p := cache.C.tail.prev
		LF:
			for p != cache.C.head {
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

func (cache threadSafeLRU) Iter(reverse bool) <-chan lruPair {
	ch := make(chan lruPair)
	go func() {
		cache.RLock()
		defer cache.RUnlock()
		if cache.C.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.C.head.next
			for p != cache.C.tail {
				ch <- lruPair{p.key, p.value}
				p = p.next
			}
		} else {
			p := cache.C.tail.prev
			for p != cache.C.head {
				ch <- lruPair{p.key, p.value}
				p = p.prev
			}
		}
		close(ch)
	}()

	return ch
}
