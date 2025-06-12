package ratelimit

import "time"

type TokenBucket struct {
	capacity    int
	fillRate    int
	tokens      int
	lastRequest time.Time
	clock       Clock
}

func NewTokenBucket(capacity, fillRate int, clock Clock) *TokenBucket {
	if clock == nil {
		clock = NewRealClock()
	}
	now := clock.Now()
	return &TokenBucket{
		capacity:    capacity,
		fillRate:    fillRate,
		tokens:      capacity,
		lastRequest: now,
		clock:       clock,
	}
}

func (tb *TokenBucket) Allow() bool {
	now := tb.clock.Now()
	elapsed := now.Sub(tb.lastRequest).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+tb.fillRate*int(elapsed))
	tb.lastRequest = now

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}
