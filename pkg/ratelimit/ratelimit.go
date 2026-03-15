package ratelimit

import (
	"context"
	"sync"
	"time"
)

type bucketEntry struct {
	tb         *TokenBucket
	lastAccess time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucketEntry
	capacity float64
	rate     float64
	ttl      time.Duration
}

func NewRateLimiter(capacity, rate float64, ttl time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucketEntry),
		capacity: capacity,
		rate:     rate,
		ttl:      ttl,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	ent, ok := rl.buckets[key]
	if !ok {
		ent = &bucketEntry{
			tb: NewBucket(
				rl.capacity,
				rl.rate,
			),
			lastAccess: time.Now(),
		}
		rl.buckets[key] = ent
	}
	ent.lastAccess = time.Now()
	tb := ent.tb
	rl.mu.Unlock()

	return tb.Allow()
}

func (rl *RateLimiter) Clean() int {
	if rl.ttl <= 0 {
		return 0
	}
	cutoff := time.Now().Add(-rl.ttl)
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n := 0
	for key, ent := range rl.buckets {
		if ent.lastAccess.Before(cutoff) {
			delete(rl.buckets, key)
			n++
		}
	}
	return n
}

func (rl *RateLimiter) StartCleanup(interval time.Duration) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				rl.Clean()
			}
		}
	}()
	return cancel
}
