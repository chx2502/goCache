package lfu

import (
	"container/list"
	"gocache/model"
	"log"
	"math"
)

type LfuCache struct {
	maxBytes int64	// 最大缓存容量
	nBytes int64	// 已使用的容量
	minFreq int64	// 记录最小使用次数
	size int
	dispatchQueue map[int64]*list.List	// 按 frequency 组织的 LRU 队列
	frequency map[string]int64	// 按 key 记录 frequency
	cache map[string]*list.Element	// <key, list.Element>
	OnEvicted func(key string, value model.Value)	// 记录被删除时的回调函数
}

type Node model.Node
type Container model.Container

func New(maxBytes int64, onEvicted func(string, model.Value)) Container {
	return &LfuCache{
		maxBytes: maxBytes,
		nBytes: 0,
		minFreq: math.MaxUint32,
		size: 0,
		dispatchQueue: make(map[int64]*list.List),
		frequency: make(map[string]int64),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (lfu *LfuCache) eliminate() {
	log.Println("eliminate()")
	for lfu.maxBytes != 0 && lfu.maxBytes < lfu.nBytes {
		lfu.Remove()
	}
}

func (lfu *LfuCache) updateMinFreq() {
	log.Println("updateMinfreq()")
	for lfu.size > 0 && lfu.dispatchQueue[lfu.minFreq].Len() == 0 {
		lfu.minFreq += 1
	}
}

func (lfu *LfuCache) removeFromDispatchQueue(key string) {
	log.Printf("removeFromDispatchQueue(%s)", key)
	freq, ok := lfu.frequency[key]
	if ok {
		lfu.dispatchQueue[freq].Remove(lfu.cache[key])
		if freq == lfu.minFreq {
			lfu.updateMinFreq()
		}
	} else {
		log.Panicf("frequency[%s] missed", key)
	}
}

func (lfu *LfuCache) addToDispatchQueue(freq int64, element *list.Element) *list.Element {
	log.Printf("removeFromDispatchQueue(%d, %v)", freq, element)
	if lfu.dispatchQueue[freq] == nil {
		lfu.dispatchQueue[freq] = list.New()
	}
	return lfu.dispatchQueue[freq].PushFront(element)
}

func (lfu *LfuCache) Get(key string) (value model.Value, ok bool) {
	log.SetPrefix("[LFU Get]")
	log.Printf("get key '%s'", key)
	if element, ok := lfu.cache[key]; ok {
		log.Printf("key '%s' hit", key)
		kv := element.Value.(*Node)
		lfu.frequency[key] += 1
		lfu.removeFromDispatchQueue(key)
		_ = lfu.addToDispatchQueue(lfu.frequency[key], element)
		return kv.Value, true
	} else {
		return nil, false
	}
}

func (lfu *LfuCache) Add(key string, value model.Value) {
	log.SetPrefix("[LFU Add]")
	log.Printf("add key '%s'", key)
	if element, ok := lfu.cache[key]; ok {
		log.Printf("key '%s' hit", key)
		kv := element.Value.(*Node)
		lfu.nBytes += int64(value.Len()) - int64(kv.Value.Len())
		if lfu.maxBytes != 0 && lfu.maxBytes < lfu.nBytes {
			lfu.eliminate()
		}
		lfu.removeFromDispatchQueue(key)
		kv.Value = value
		lfu.frequency[key] += 1
		_ = lfu.addToDispatchQueue(lfu.frequency[key], element)
		log.Printf("updated new key '%s'", key)
	} else {
		log.Printf("key '%s' missed", key)
		lfu.nBytes += int64(len(key)) + int64(value.Len())
		if lfu.maxBytes != 0 && lfu.maxBytes < lfu.nBytes {
			lfu.eliminate()
		}
		var freq int64 = 1
		lfu.frequency[key] = freq
		if lfu.dispatchQueue[freq] == nil {
			lfu.dispatchQueue[freq] = list.New()
		}
		element := lfu.dispatchQueue[freq].PushFront(&Node{Key: key, Value: value})
		lfu.cache[key] = element
		if freq < lfu.minFreq {
			lfu.minFreq = freq
		}
		lfu.size += 1
		log.Printf("added new key '%s'", key)
	}
}

func (lfu *LfuCache) Remove() {
	log.Println("LFU remove")
	element := lfu.dispatchQueue[lfu.minFreq].Back()
	if element != nil {
		lfu.dispatchQueue[lfu.minFreq].Remove(element)
		kv := element.Value.(*Node)
		delete(lfu.cache, kv.Key)
		delete(lfu.frequency, kv.Key)
		lfu.nBytes -= int64(len(kv.Key)) + int64(kv.Value.Len())
		if lfu.OnEvicted != nil {
			lfu.OnEvicted(kv.Key, kv.Value)
		}
		lfu.size -= 1
		lfu.updateMinFreq()
		log.Printf("key '%s' removed", kv.Key)
	}
}

func (lfu *LfuCache) Len() int {
	return lfu.size
}