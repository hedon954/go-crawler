package engine

import (
	"sync"

	"github.com/hedon954/go-crawler/fetcher"
	"go.uber.org/zap"
)

// Crawler represents the global crawl instance
type Crawler struct {
	out chan fetcher.ParseResult

	// store the visited fetcher.Request
	Visited     map[string]bool
	VisitedLock sync.Mutex

	options
}

func NewCrawler(opts ...Option) *Crawler {
	dopts := defaultOptions
	for _, opt := range opts {
		opt(&dopts)
	}
	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	e.out = make(chan fetcher.ParseResult)
	e.options = dopts
	return e
}

func (c *Crawler) HasVisited(r *fetcher.Request) bool {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()
	return c.Visited[r.UniqueSign()]
}

func (c *Crawler) StoreVisited(reqs ...*fetcher.Request) {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()

	for _, r := range reqs {
		c.Visited[r.UniqueSign()] = true
	}
}

func (c *Crawler) Run() {
	go c.Schedule()
	for i := 0; i < c.WorkCount; i++ {
		go c.CreateWork()
	}
	c.HandleResult()
}

func (c *Crawler) Schedule() {
	var reqs []*fetcher.Request
	for _, seed := range c.Seeds {
		seed.RootReq.Task = seed
		seed.RootReq.Url = seed.Url
		reqs = append(reqs, seed.RootReq)
	}
	go c.scheduler.Schedule()
	go c.scheduler.Push(reqs...)
}

func (c *Crawler) CreateWork() {
	for {
		r := c.scheduler.Pull()
		if err := r.Check(); err != nil {
			c.Logger.Error("check failed", zap.Error(err))
			continue
		}

		// Remove duplicate request
		if c.HasVisited(r) {
			continue
		}
		c.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if err != nil {
			c.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			continue
		}
		if len(body) < 6000 {
			c.Logger.Error("can't fetch ",
				zap.Int("length", len(body)),
				zap.String("url", r.Url),
			)
			continue
		}

		result := r.ParseFunc(body, r)
		if len(result.Requests) > 0 {
			go c.scheduler.Push(result.Requests...)
		}
		c.out <- result
	}
}

func (c *Crawler) HandleResult() {
	for {
		select {
		case result := <-c.out:
			for _, item := range result.Items {
				c.Logger.Sugar().Info("get result:", item)
			}
		}
	}
}
