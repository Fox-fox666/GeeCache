package RWMutexCache

import (
	"GeeCache/LruCache"
	"sync"
)

type RWMutexCache struct {
	rwm      sync.Mutex
	cache    *LruCache.Lru_Cache
	MaxBytes int64
}

func (c *RWMutexCache) Get(key string) (v BytesData, ok bool) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	if c.cache == nil {
		return
	}
	var data LruCache.Cachedata
	if data, ok = c.cache.Get(key); ok {
		return data.(BytesData), ok
	}
	return
}

func (c *RWMutexCache) Add(key string, value BytesData) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	if c.cache == nil {
		c.cache = LruCache.NewCache(c.MaxBytes, nil) //延迟初始化，当第一次要用到缓存的时候才创建缓存
	}

	c.cache.Add(key, value)
	return
}
