package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu          *sync.RWMutex
	CachedEntry map[string]cacheEntry
}

func (c *Cache) reapLoop(interval float64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.CachedEntry {
			if time.Now().Sub(v.createdAt).Seconds() > interval {
				delete(c.CachedEntry, k)
			}
		}
		c.mu.Unlock()
	}
	defer ticker.Stop()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.CachedEntry[key]
	if !ok {
		return nil, false
	}
	return v.val, true
}

func (c *Cache) Set(key string, val []byte) {
	c.mu.Lock()
	c.CachedEntry[key] = cacheEntry{time.Now(), val}
	c.mu.Unlock()
	return
}

func NewCache(interval float64) *Cache {
	newCache := &Cache{
		mu:          &sync.RWMutex{},
		CachedEntry: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	return newCache
}
