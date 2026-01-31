package ratelimit

import (
	"errors"
	"math"
	"sync"
	"time"
)

type RaterLimiter interface {
	Deduct(id string, cost float64) (float64, error)
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

func (b *Bucket) consume(cost, maxValue float64, refillIntervalMillis int64) (float64, bool) {
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
		return b.count, true
	}
	return b.count, false
}

type TokenBucketRateLimit struct {
	refillIntervalMillis int64
	maxValue             float64
	buckets              map[string]*Bucket
	mu                   sync.RWMutex
}

func (tbr *TokenBucketRateLimit) Deduct(id string, cost float64) (float64, error) {
	now := time.Now().UTC().UnixMilli()

	// Check if bucket exists (read lock)
	tbr.mu.RLock()
	bucket, ok := tbr.buckets[id]
	tbr.mu.RUnlock()

	if !ok {
		// Create new bucket (write lock)
		tbr.mu.Lock()
		// Double-check after acquiring write lock (another goroutine might have created it)
		bucket, ok = tbr.buckets[id]
		if !ok {
			newBucket := &Bucket{count: tbr.maxValue, refilledAt: now}
			tbr.buckets[id] = newBucket
		}
		tbr.mu.Unlock()
		return tbr.Deduct(id, cost)
	}

	remaining, ok := bucket.consume(cost, tbr.maxValue, tbr.refillIntervalMillis)
	if !ok {
		return remaining, errors.New("rate limit reached")
	}
	return remaining, nil
}
