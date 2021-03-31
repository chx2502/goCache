package lru

import "container/list"

type Cache struct {
	maxBytes int64	// 最大缓存容量
	nBytes int64	// 已使用的容量
	linkedList *list.List
	cache map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type node struct {
	key   string
	value Value
}

type Value interface {
	Len() int	// 返回 value 占用内存的大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		nBytes: 0,
		linkedList: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(element)
		kv := element.Value.(*node) // list.Element存储的是任意类型，需要进行类型断言
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	element := c.linkedList.Back()
	if element != nil {
		c.linkedList.Remove(element)
		kv := element.Value.(*node)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(element)
		kv := element.Value.(*node)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := c.linkedList.PushFront(&node{key, value})
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.linkedList.Len()
}