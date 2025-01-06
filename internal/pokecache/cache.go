package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu    *sync.Mutex
	store map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		store: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value := cacheEntry{}
	value, exists := c.store[key]
	if exists {
		return value.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.store {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.store, k)
		}
	}
}
