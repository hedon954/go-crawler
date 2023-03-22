package fetcher

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/hedon954/go-crawler/collector"
	"github.com/hedon954/go-crawler/limiter"
)

// Task represents a complete crawl task
type Task struct {
	Property

	Visited     map[string]bool
	VisitedLock sync.Mutex

	RootReq *Request
	Fetcher Fetcher
	Rule    RuleTree

	Logger  *zap.Logger
	Storage collector.Store
	Limiter limiter.MultiLimiter
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
	// Headers needs to be added to http headers
	Headers map[string]string `json:"headers"`
}
