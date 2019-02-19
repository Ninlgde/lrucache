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
			//expires := cache.cap >> 2

			// 思考：如果淘汰的时候淘汰1/4。。虽然淘汰的次数少了，但是淘汰的那一次会从O(1)的操作增加到O(n/4)
			// 如果n特别特别大，那么这一次的操作会非常耗时。平均是O(2)
			// 如果改回满了，删一次，那么满了之后的增加多了一次删除操作，大约是O(2)
			// 综上，还是选择每次删除一个。
			expires := 1
			//for i := 0; i < expires; i++ {
			rmnode := cache.head.next
			rmkey := rmnode.key       // 找到对应的key
			delete(cache.dict, rmkey) // 一定要把map里的key给删除

			freeNode(rmnode)
			//}
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
	node, ok := cache.dict[k]
	if ok {
		// 命中，将此node移到双向链表的末尾
		// 1.先从原位置删除
		value := freeNode(node)
		// add 的时候有find操作，会改变value
		if v != nil {
			value = v
		}

		// 2.再添加到尾部
		cache.add(k, value, true) // 因为命中了 所以size不自增
		return value
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

	value := freeNode(node)

	cache.len--
	return value
}

/**
还不太清楚go的垃圾回收机制
理论上要把node所有引用的地方都制空才会被回收吧
*/
func freeNode(node *lruNode) lruValue {
	// 把指针操作也放到里面
	node.next.prev = node.prev
	node.prev.next = node.next
	v := node.value
	// 把node的引用也置空
	// 其实没有必要，golang的回收是检查对象是否被引用，上面的操作已经完成了解引用
	// 所以下面理论上不需要，但还是加上吧
	node.value = nil
	node.key = nil
	node.next = nil
	node.prev = nil
	_ = node
	return v
}

// 经过下面的方法测试，在不添加freeNode时，发生了内存泄漏，添加后内存泄漏问题消失
//func main() {
//
//	cache := lru.NewLRUCache(100)
//
//	i:=0
//	for {
//		cache.Add(i,i)
//		i++
//	}
//}
// 后记：内存泄漏时由于添加的bug导致的，跟freeNode没关系
// 需要好好看看golang的内存回收了！！！
// 后记2：
// 		golang gc： http://legendtkl.com/2017/04/28/golang-gc/
