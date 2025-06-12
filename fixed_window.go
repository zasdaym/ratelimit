package ratelimit

import "time"

type FixedWindow struct {
	windowSize      int
	maxRequests     int
	requests        int
	lastWindowReset time.Time
	clock           Clock
}

func NewFixedWindow(windowSize, maxRequests int) *FixedWindow {
	clock := NewRealClock()
	return &FixedWindow{
		windowSize:      windowSize,
		maxRequests:     maxRequests,
		lastWindowReset: clock.Now(),
		clock:           clock,
	}
}

func (fw *FixedWindow) Allow() bool {
	now := fw.clock.Now()
	if int(now.Sub(fw.lastWindowReset).Seconds()) > fw.windowSize {
		fw.requests = 0
		fw.lastWindowReset = now
	}

	if fw.requests < fw.maxRequests {
		fw.requests++
		return true
	}
	return false
}
