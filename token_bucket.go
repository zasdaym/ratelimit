package ratelimit

import "time"

type TokenBucket struct {
	capacity  int
	fillRate  int
	tokens    int
	lastAllow time.Time
	clock     Clock
}

func NewTokenBucket(capacity, fillRate int) *TokenBucket {
	clock := NewRealClock()
	now := clock.Now()
	return &TokenBucket{
		capacity:  capacity,
		fillRate:  fillRate,
		tokens:    capacity,
		lastAllow: now,
		clock:     clock,
	}
}

func (tb *TokenBucket) Allow() bool {
	now := tb.clock.Now()
	elapsed := now.Sub(tb.lastAllow).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+tb.fillRate*int(elapsed))

	if tb.tokens > 0 {
		tb.tokens--
		tb.lastAllow = now
		return true
	}
	return false
}
