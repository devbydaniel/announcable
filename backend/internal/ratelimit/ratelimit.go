package ratelimit

import (
	"errors"
	"math"
	"sync"
	"time"
)

type RaterLimiter interface {
	Deduct(id string, cost float64) error
}

func New(refillIntervalSeconds int64, maxValue float64) RaterLimiter {
	buckets := make(map[string]*Bucket)
	return &TokenBucketRateLimit{refillIntervalMillis: refillIntervalSeconds * 1000, maxValue: maxValue, buckets: buckets}
}

type Bucket struct {
	count      float64
	refilledAt int64
	mu         sync.Mutex
}

func (b *Bucket) consume(cost, maxValue float64, refillIntervalMillis int64) bool {
	// lock the bucket
	b.mu.Lock()
	defer b.mu.Unlock()
	// calculate refill until now & refill bucket
	now := time.Now().UTC().UnixMilli()
	refill := float64(now-b.refilledAt) / float64(refillIntervalMillis) * maxValue
	b.count = math.Min(b.count+refill, maxValue)
	b.refilledAt = now

	// deduct cost if possible
	if b.count >= cost {
		b.count -= cost
		return true
	}
	return false
}

type TokenBucketRateLimit struct {
	refillIntervalMillis int64
	maxValue             float64
	buckets              map[string]*Bucket
}

func (tbr *TokenBucketRateLimit) Deduct(id string, cost float64) error {
	now := time.Now().UTC().UnixMilli()
	// check if userId is already part of the map
	bucket, ok := tbr.buckets[id]
	if !ok {
		bucket := Bucket{count: tbr.maxValue, refilledAt: now}
		tbr.buckets[id] = &bucket
		return tbr.Deduct(id, cost)
	}
	if ok := bucket.consume(cost, tbr.maxValue, tbr.refillIntervalMillis); !ok {
		return errors.New("rate limit reached")
	}
	return nil
}
