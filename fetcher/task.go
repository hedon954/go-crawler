package fetcher

import (
	"github.com/hedon954/go-crawler/collector"
	"sync"
	"time"
)

// Task represents a complete crawl task
type Task struct {

	// The unique signature of the Task
	Name     string
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
	Rule    RuleTree

	Storage collector.Store
}
