package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	refillRate float64
	last       time.Time
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		last:       time.Now(),
	}
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	diff := now.Sub(t.last).Seconds()
	if diff > 0 {
		t.tokens += diff * t.refillRate
		t.last = now
		if t.tokens > t.capacity {
			t.tokens = t.capacity
		}
	}
	if t.tokens >= 1 {
		t.tokens--
		return true
	} else {
		return false
	}
}
