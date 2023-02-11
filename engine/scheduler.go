package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
)

// Scheduler defines the behavior of scheduing crawl request
type Scheduler interface {

	// Schedule starts the scheduler
	Schedule()

	// Push the request into scheduler
	Push(...*fetcher.Request)

	// Pull a request from scheduler
	Pull() *fetcher.Request
}
