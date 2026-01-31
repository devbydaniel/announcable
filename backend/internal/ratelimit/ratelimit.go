package ratelimit

import (
	"errors"
	"math"
	"sync"
	"time"
)

// RaterLimiter defines the interface for a rate limiter.
type RaterLimiter interface {
	Deduct(id string, cost float64) error
}

// New creates a new token bucket rate limiter with the given refill interval and max value.
func New(refillIntervalSeconds int64, maxValue float64) RaterLimiter {
	buckets := make(map[string]*Bucket)
	return &TokenBucketRateLimit{refillIntervalMillis: refillIntervalSeconds * 1000, maxValue: maxValue, buckets: buckets}
}

// Bucket holds the token count and refill state for a single rate-limited entity.
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

// TokenBucketRateLimit implements RaterLimiter using the token bucket algorithm.
type TokenBucketRateLimit struct {
	refillIntervalMillis int64
	maxValue             float64
	buckets              map[string]*Bucket
}

// Deduct attempts to consume tokens from the bucket for the given ID, returning an error if the limit is exceeded.
func (tbr *TokenBucketRateLimit) Deduct(id string, cost float64) error {
	now := time.Now().UTC().UnixMilli()
	// check if userID is already part of the map
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
