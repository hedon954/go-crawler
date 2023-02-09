package fetcher

import (
	"time"
)

// Task represents a complete crawl task
type Task struct {
	Url      string
	Cookie   string
	WaitTime time.Duration
	MaxDepth int
	RootReq  *Request
	Fetcher  Fetcher
}
