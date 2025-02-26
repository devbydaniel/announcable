package memcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cacher interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiresIn time.Duration)
	Flush()
}

func New(defaultExpiration, cleanupInterval time.Duration) Cacher {
	c := cache.New(defaultExpiration, cleanupInterval)
	return c
}
