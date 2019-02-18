package lru

/**
lru node
双向列表的节点
*/
type lruNode struct {
	next  *lruNode // 后指针
	prev  *lruNode // 前指针
	value lruValue // 缓存的值
	key   lruKey   // 缓存的key
}

/**
LRU 缓存
由一个map和一个双向链表组成
可将查找、添加等操作的时间复杂度较少到O(1) (理论上，取决于map的实现)
*/
type threadUnsafeLRU struct {
	head *lruNode            // 头指针
	tail *lruNode            // 尾指针
	dict map[lruKey]*lruNode // 存放数据的 map，提高查找效率
	len  int                 // 当前数量
	cap  int                 // 总量
}

func newThreadUnsafeLRU() *threadUnsafeLRU {
	return &threadUnsafeLRU{}
}

/**
创建缓存
cap: 容量，缓存最多存多少数据
*/
func (cache *threadUnsafeLRU) Create(cap int) {
	// 初始化尾指针
	cache.tail = &lruNode{
		next:  nil,
		prev:  nil,
		value: nil,
		key:   nil,
	}
	// 初始化头指针，要将头的后指针指向尾指针
	cache.head = &lruNode{
		next:  cache.head,
		prev:  nil,
		value: nil,
		key:   nil,
	}
	cache.tail.prev = cache.head           // 完成头尾相连
	cache.dict = make(map[lruKey]*lruNode) // init map
	cache.len = 0
	cache.cap = cap

}

/**
添加一个元素
k: key
v: value

cost: O(1)(base on map's implement)
*/
func (cache *threadUnsafeLRU) Add(k lruKey, v lruValue) {
	// 先查找是否在缓存里，如果有就不添加了
	if cache.find(k, v) != nil {
		return
	}
	// add
	cache.add(k, v, false)
}

/**
查找一个元素
k: key
return: value or nil

cost: O(1)(base on map's implement)
*/
func (cache *threadUnsafeLRU) Find(k lruKey) lruValue {
	return cache.find(k, nil)
}

/**
当前缓存大小
return: size of cache

cost: O(1)
*/
func (cache *threadUnsafeLRU) Size() int {
	return cache.len
}

/**
删除一个元素
k: key
return: value or nil
*/
func (cache *threadUnsafeLRU) Remove(k lruKey) lruValue {
	node := cache.find(k, nil) // this step will move k to tail
	if node == nil {
		return nil
	}
	// k in cache
	return cache.poptail(k)
}

/**
遍历缓存中所有的数据的迭代器
reverse: 是否翻转 true = 正序 false = 倒序(默认，淘汰的是从头部，所以从后往前是默认)
return: 迭代器 func
*/
func (cache *threadUnsafeLRU) Iterator(reverse bool) *Iterator {
	iterator, ch, stopCh := newIterator(cache.cap)
	go func() {
		if cache.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.head.next
		LT:
			for p != cache.tail {
				select {
				case <-stopCh:
					break LT
				case ch <- lruPair{p.key, p.value}:
				}
				p = p.next
			}
		} else {
			p := cache.tail.prev
		LF:
			for p != cache.head {
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

func (cache *threadUnsafeLRU) Iter(reverse bool) <-chan lruPair {
	ch := make(chan lruPair, cache.cap)
	go func() {
		if cache.len == 0 {
			close(ch)
			return
		}
		if reverse {
			p := cache.head.next
			for p != cache.tail {
				ch <- lruPair{p.key, p.value}
				p = p.next
			}
		} else {
			p := cache.tail.prev
			for p != cache.head {
				ch <- lruPair{p.key, p.value}
				p = p.prev
			}
		}
		close(ch)
	}()

	return ch
}

/**
内置的add方法
cache: 缓存
k: key
v: value
notinc: 是否增加 如果命中的话，只是移动位置，不需要增加size
*/
func (cache *threadUnsafeLRU) add(k lruKey, v lruValue, notinc bool) {
	if !notinc {
		cache.len += 1
		if cache.len > cache.cap {
			// 为提高效率，每次从头部淘汰整体的1/4
			expires := cache.cap >> 2
			for i := 0; i < expires; i++ {
				rmkey := cache.head.next.key // 找到对应的key
				delete(cache.dict, rmkey)    // 一定要把map里的key给删除

				// head指针后移一位
				cache.head.next = cache.head.next.next
				cache.head.prev = cache.head
			}
			cache.len -= expires // size减小到删除后的真实size
		}
	}
	// 创建node，并添加到尾部和map中
	node := &lruNode{
		next:  cache.tail,
		prev:  cache.tail.prev,
		value: v,
		key:   k,
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
func (cache *threadUnsafeLRU) find(k lruKey, v lruValue) lruValue {
	old, ok := cache.dict[k]
	if ok {
		// 命中，将此node移到双向链表的末尾
		// 1.先从原位置删除
		node := old
		node.prev.next = node.next // 前的后是后
		node.next.prev = node.prev // 后的前是前

		// add 的时候有find操作，会改变value
		new := node.value
		if v != nil {
			new = v
		}

		// 2.再添加到尾部
		cache.add(k, new, true) // 因为命中了 所以size不自增
		return node.value
	}
	return nil // 没有命中返回nil
}

/**
从尾部pop出一个node
cache: 缓存
k: key
return: pop的value or nil(空表时)
*/
func (cache *threadUnsafeLRU) poptail(k lruKey) lruValue {
	if cache.len == 0 {
		return nil
	}
	node := cache.tail.prev
	rmkey := node.key
	delete(cache.dict, rmkey)
	// tail 指针向前移动一位
	cache.tail.prev = node.prev
	node.prev.next = cache.tail

	cache.len--
	return node.value
}
