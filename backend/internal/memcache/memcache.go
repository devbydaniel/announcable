package memcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cacher defines the interface for an in-memory key-value cache.
type Cacher interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiresIn time.Duration)
	Flush()
}

// New creates a new in-memory cache with the given expiration and cleanup interval.
func New(defaultExpiration, cleanupInterval time.Duration) Cacher {
	c := cache.New(defaultExpiration, cleanupInterval)
	return c
}
