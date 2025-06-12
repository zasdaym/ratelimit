package ratelimit

import (
	"testing"
	"time"
)

func TestFixedWindowAllow(t *testing.T) {
	t.Parallel()

	var (
		windowSize  = 2
		maxRequests = 5
		clock       = NewMockClock(time.Now())
		bucket      = NewFixedWindow(windowSize, maxRequests, clock)
	)

	// Test initial window
	for i := 0; i < maxRequests; i++ {
		if !bucket.Allow() {
			t.Errorf("request %d should be allowed", i+1)
		}
	}
	if bucket.Allow() {
		t.Error("request should be denied when max requests are reached")
	}

	// Test window reset
	clock.Advance(time.Duration(windowSize+1) * time.Second)
	if !bucket.Allow() {
		t.Error("request should be allowed after window reset")
	}

	// Test partial window advance
	clock.Advance(time.Duration(windowSize/2) * time.Second)
	for i := 0; i < maxRequests-1; i++ {
		if !bucket.Allow() {
			t.Errorf("request %d should be allowed in new window", i+1)
		}
	}
	if bucket.Allow() {
		t.Error("request should be denied when max requests are reached in new window")
	}
}
