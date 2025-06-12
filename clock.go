package ratelimit

import "time"

// Clock provides current time.
type Clock interface {
	Now() time.Time
}

// RealClock is the real implementation of Clock.
type RealClock struct{}

func NewRealClock() *RealClock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

// FakeClock is a fake implementation of Clock for testing.
type FakeClock struct {
	currentTime time.Time
}

func NewFakeClock(initialTime time.Time) *FakeClock {
	return &FakeClock{
		currentTime: initialTime,
	}
}

func (c *FakeClock) Now() time.Time {
	return c.currentTime
}

// Advance the time by the given duration.
func (c *FakeClock) Advance(d time.Duration) {
	c.currentTime = c.currentTime.Add(d)
}
