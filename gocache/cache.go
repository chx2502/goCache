package gocache

import (
	"gocache/lfu"
	"gocache/model"
	"gocache/lru"
	"strings"
	"sync"
)

type cache struct {
	mu sync.Mutex
	strategy string
	container model.Container
	cacheBytes int64
}
type Value model.Value
func NewContainer(maxBytes int64, strategy string, onEvicted func(string, model.Value)) (container model.Container) {
	if strings.Compare(strategy, "lru") == 0 {
		container = lru.New(maxBytes, onEvicted)
	} else if strings.Compare(strategy, "lfu") == 0 {
		container = lfu.New(maxBytes, onEvicted)
	}
	return container
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.container == nil {
		c.container = NewContainer(c.cacheBytes, c.strategy, nil)
	}
	c.container.Add(key, value)

}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.container == nil {
		return
	}

	if v, ok := c.container.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

