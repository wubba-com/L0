package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCache(defaultExpiration, cleanupInterval time.Duration) Cache {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &localCache{cache: c}
}

var NoExpiration = cache.NoExpiration

type Cache interface {
	Set(string, interface{}, time.Duration)
	Get(string) (interface{}, bool)
}

type localCache struct {
	cache *cache.Cache
}

func (c *localCache) Set(key string, value interface{}, duration time.Duration) {
	c.cache.Set(key, value, duration)
}

func (c *localCache) Get(key string) (interface{}, bool) {
	v, _ := c.cache.Get(key)
	if v == nil {
		return nil, false
	}

	return v, true
}