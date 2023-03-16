package limiter

import (
	"time"

	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

// Per calculates the rate of the limiter according to eventCount and duration
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
