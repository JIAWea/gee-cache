package lru

import "container/list"

// Cache 实现 LRU 淘汰算法
// 维护一个双向链表，每当访问过一个 key，将把该 key 对应的节点移到队尾
type Cache struct {
	maxBytes int64      // 允许最大内存
	nBytes   int64      // 当前已用内存
	ll       *list.List // 双向链表
	cache    map[string]*list.Element
	// OnEvicted 回调
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value 计算需要多少字节
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 获取值
// 1. 从字典中获取对应的双向链表的节点
// 2. 将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰，移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	// 取出队首
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 更新当前所用内存
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 判断并清理内存
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len 查询有多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
