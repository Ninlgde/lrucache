package lru

import (
	"sync"
)

// 值的类型 类似void*
type lruvalue interface {}
// key的类型  类似void*
type lrukey interface {}

/**
	lru node
	双向列表的节点
 */
type lrunode struct {
	next *lrunode						// 后指针
	prev *lrunode						// 前指针
	value lruvalue						// 缓存的值
	key lrukey							// 缓存的key
}

/**
	LRU 缓存
	由一个map和一个双向链表组成
	可将查找、添加等操作的时间复杂度较少到O(1) (理论上，取决于map的实现)
 */
type LRUCache struct {
	head *lrunode 						// 头指针
	tail *lrunode						// 尾指针
	dict map[lrukey] *lrunode			// 存放数据的map，提高查找效率
	len int								// 当前数量
	cap int								// 总量
	mutex *sync.RWMutex					// 协程锁
}

/**
	创建缓存
	cap: 容量，缓存最多存多少数据
 */
func (cache *LRUCache) Create(cap int) {
	// 初始化尾指针
	cache.tail = &lrunode{
		next: nil,
		prev: nil,
		value: nil,
		key: nil,
	}
	// 初始化头指针，要将头的后指针指向尾指针
	cache.head = &lrunode{
		next: cache.head,
		prev: nil,
		value: nil,
		key: nil,
	}
	cache.tail.prev = cache.head // 完成头尾相连
	cache.dict = make(map[lrukey] *lrunode) // init map
	cache.len = 0
	cache.cap = cap

	cache.mutex = &sync.RWMutex{} // init lock
}

/**
	添加一个元素
	k: key
	v: value

	cost: O(1)(base on map's implement)
 */
func (cache *LRUCache) Add(k lrukey, v lruvalue) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	// 先查找是否在缓存里，如果有就不添加了
	if _find(cache, k, v) != nil {
		return
	}
	// add
	_add(cache, k, v, false)
}

/**
	查找一个元素
	k: key
	return: value or nil

	cost: O(1)(base on map's implement)
 */
func (cache *LRUCache) Find(k lrukey) lruvalue {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	return _find(cache, k, nil)
}

/**
	当前缓存大小
	return: size of cache

	cost: O(1)
 */
func (cache *LRUCache) Size() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return cache.len
}

/**
	遍历缓存中所有的数据
	不清楚是否需要加锁
	reverse: 是否翻转 true = 正序 false = 倒序(默认，淘汰的是从头部，所以从后往前是默认)
	return: 迭代器 func
 */
func (cache LRUCache) Iterator(reverse bool) func() lruvalue {
	var p *lrunode
	if reverse {
		p = cache.head
	} else {
		p = cache.tail
	}
	return func() lruvalue {
		if p != nil {
			if reverse {
				p = p.next
			} else {
				p = p.prev
			}
			return p.value
		}
		return nil
	}
}

/**
	内置的add方法
	cache: 缓存
	k: key
	v: value
	notinc: 是否增加 如果命中的话，只是移动位置，不需要增加size
 */
func _add(cache *LRUCache, k lrukey, v lruvalue, notinc bool) {
	if !notinc {
		cache.len += 1
		if cache.len > cache.cap {
			// 为提高效率，每次从头部淘汰整体的1/4
			expires := cache.cap << 2
			for i:=0; i<expires; i++ {
				movekey := cache.head.key // 找到对应的key
				delete(cache.dict, movekey) // 一定要把map里的key给删除

				// head指针后移一位
				cache.head.next = cache.head.next.next
				cache.head.prev = cache.head
			}
			cache.len -= expires // size减小到删除后的真实size
		}
	}
	// 创建node，并添加到尾部和map中
	node := &lrunode{
		next: cache.tail,
		prev: cache.tail.prev,
		value: v,
		key: k,
	}

	cache.tail.prev.next = node
	cache.tail.prev = node
	cache.dict[k] = node
}

/**
	内置查找方法
	cache: 缓存
	k: key
	v: value add有find操作，有可能改变value
	return: 找到的value or nil
 */
func _find(cache *LRUCache, k lrukey, v lruvalue) lruvalue {
	old, ok := cache.dict[k]
	if ok {
		// 命中，将此node移到双向链表的末尾
		node := old
		node.prev.next = node.next // 前的后是后
		node.next.prev = node.prev // 后的前是前

		// add 的时候有find操作，会改变value
		new := node.value
		if v != nil {
			new = v
		}

		_add(cache, k, new, true) // 因为命中了 所以size不自增
		return node.value
	}
	return nil // 没有命中返回nil
}
