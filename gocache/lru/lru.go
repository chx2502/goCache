package lru

import (
	"container/list"
	"gocache/model"
)

type Node model.Node
type Container model.Container

type LruCache struct {
	maxBytes int64	// 最大缓存容量
	nBytes int64	// 已使用的容量
	linkedList *list.List	// 维护 LRU 队列
	cache map[string]*list.Element	// key-listNode
	OnEvicted func(key string, value model.Value)	// 记录被删除时的回调函数
}

func New(maxBytes int64, onEvicted func(string, model.Value)) Container {
	container := &LruCache{
		maxBytes: maxBytes,
		nBytes: 0,
		linkedList: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
	return container
}

func (lru *LruCache) Get(key string) (value model.Value, ok bool) {
	if element, ok := lru.cache[key]; ok {
		lru.linkedList.MoveToFront(element)
		kv := element.Value.(*Node) // list.Element存储的是任意类型，需要进行类型断言
		return kv.Value, true
	}
	return nil, false
}

func (lru *LruCache) Remove() {
	element := lru.linkedList.Back()
	if element != nil {
		lru.linkedList.Remove(element)
		kv := element.Value.(*Node)
		delete(lru.cache, kv.Key)
		lru.nBytes -= int64(len(kv.Key)) + int64(kv.Value.Len())
		if lru.OnEvicted != nil {
			lru.OnEvicted(kv.Key, kv.Value)
		}
	}
}

func (lru *LruCache) Add(key string, value model.Value) {
	if element, ok := lru.cache[key]; ok {
		lru.linkedList.MoveToFront(element)
		kv := element.Value.(*Node)
		lru.nBytes += int64(value.Len()) - int64(kv.Value.Len())
		kv.Value = value
	} else {
		element := lru.linkedList.PushFront(&Node{key, value})
		lru.cache[key] = element
		lru.nBytes += int64(len(key)) + int64(value.Len())
	}
	for lru.maxBytes != 0 && lru.maxBytes < lru.nBytes {
		lru.Remove()
	}
}

func (lru *LruCache) Len() int {
	return lru.linkedList.Len()
}