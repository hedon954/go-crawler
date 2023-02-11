package fetcher

import (
	"sync"
	"time"
)

// Task represents a complete crawl task
type Task struct {
	Url      string
	Cookie   string
	WaitTime time.Duration
	MaxDepth int

	// Mark whether the site can be crawled repeated
	Reload bool

	Visited     map[string]bool
	VisitedLock sync.Mutex

	RootReq *Request
	Fetcher Fetcher
}
