package limiter

import (
	"context"
	"sort"

	"golang.org/x/time/rate"
)

type MultiLimiter struct {
	limiters []RateLimiter
}

func NewMultiLimiter(limiters ...RateLimiter) *MultiLimiter {
	sort.Slice(limiters, func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	})
	return &MultiLimiter{limiters: limiters}
}

func (ml *MultiLimiter) Wait(ctx context.Context) error {
	for _, l := range ml.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ml *MultiLimiter) Limit() rate.Limit {
	if len(ml.limiters) < 1 {
		return 0
	}
	return ml.limiters[0].Limit()
}
