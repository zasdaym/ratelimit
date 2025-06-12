package ratelimit

import (
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	t.Parallel()

	var (
		capacity = 10
		fillRate = 5
		bucket   = NewTokenBucket(capacity, fillRate)
	)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	if bucket.tokens != capacity {
		t.Errorf("expected initial tokens %d, got %d", capacity, bucket.tokens)
	}
}

func TestTokenBucketAllowInitialTokens(t *testing.T) {
	t.Parallel()

	var (
		capacity = 3
		fillRate = 1
		bucket   = NewTokenBucket(capacity, fillRate)
	)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	for i := 0; i < capacity; i++ {
		if !bucket.Allow() {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	if bucket.Allow() {
		t.Error("request should be denied when bucket is empty")
	}
}

func TestTokenBucketRefill(t *testing.T) {
	t.Parallel()

	bucket := NewTokenBucket(2, 1)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	bucket.Allow()
	bucket.Allow()

	if bucket.Allow() {
		t.Error("request should be denied when bucket is empty")
	}

	clock.Advance(2 * time.Second)
	if !bucket.Allow() {
		t.Error("request should be allowed after refill")
	}
}

func TestTokenBucketCapacityLimit(t *testing.T) {
	t.Parallel()

	bucket := NewTokenBucket(2, 10)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	bucket.Allow()

	clock.Advance(2 * time.Second)

	allowedRequests := 0
	for i := 0; i < 5; i++ {
		if bucket.Allow() {
			allowedRequests++
		}
	}

	if allowedRequests != 2 {
		t.Errorf("expected 2 allowed requests due to capacity limit, got %d", allowedRequests)
	}
}

func TestTokenBucketZeroCapacity(t *testing.T) {
	t.Parallel()

	bucket := NewTokenBucket(0, 1)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	if bucket.Allow() {
		t.Error("request should be denied with zero capacity")
	}
}

func TestTokenBucketZeroFillRate(t *testing.T) {
	t.Parallel()

	bucket := NewTokenBucket(1, 0)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	if !bucket.Allow() {
		t.Error("first request should be allowed")
	}

	clock.Advance(time.Second)
	if bucket.Allow() {
		t.Error("second request should be denied with zero fill rate")
	}
}

func TestTokenBucketConcurrentAccess(t *testing.T) {
	t.Parallel()

	bucket := NewTokenBucket(100, 10)

	clock := NewFakeClock(time.Now())
	bucket.clock = clock

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 50; i++ {
			bucket.Allow()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 50; i++ {
			bucket.Allow()
		}
		done <- true
	}()

	<-done
	<-done
}
