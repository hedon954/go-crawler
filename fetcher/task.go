package fetcher

import (
	"github.com/hedon954/go-crawler/collector"
	"sync"
	"time"
)

// Task represents a complete crawl task
type Task struct {
	Property

	Visited     map[string]bool
	VisitedLock sync.Mutex

	RootReq *Request
	Fetcher Fetcher
	Rule    RuleTree

	Storage collector.Store
}

type Property struct {
	// The unique signature of the Task
	Name     string        `json:"name"`
	Url      string        `json:"url"`
	Cookie   string        `json:"cookie"`
	WaitTime time.Duration `json:"wait_time"`
	// Mark whether the site can be crawled repeated
	Reload   bool  `json:"reload"`
	MaxDepth int64 `json:"max_depth"`
}
