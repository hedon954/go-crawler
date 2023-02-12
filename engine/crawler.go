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

	// try again when crawled failed
	failures    map[string]*fetcher.Request
	failureLock sync.Mutex

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
	e.failures = make(map[string]*fetcher.Request)
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
		task := Store.hash[seed.Name]
		task.Fetcher = seed.Fetcher
		rootReqs := task.Rule.Root()
		for _, req := range rootReqs {
			req.Task = task
		}
		reqs = append(reqs, rootReqs...)
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
		if !r.Task.Reload && c.HasVisited(r) {
			continue
		}
		c.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if err != nil {
			c.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			c.SetFailure(r)
			continue
		}
		if len(body) < 6000 {
			c.Logger.Error("can't fetch ",
				zap.Int("length", len(body)),
				zap.String("url", r.Url),
			)
			c.SetFailure(r)
			continue
		}

		// Get the rule corresponding to the request
		rule := r.Task.Rule.Trunk[r.RuleName]
		result := rule.ParseFunc(&fetcher.Context{
			Body: body,
			Req:  r,
		})
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

func (c *Crawler) SetFailure(req *fetcher.Request) {
	if !req.Task.Reload {
		c.VisitedLock.Lock()
		delete(c.Visited, req.UniqueSign())
		c.VisitedLock.Unlock()
	}
	c.failureLock.Lock()
	defer c.failureLock.Unlock()
	// retry at first failure
	if _, ok := c.failures[req.UniqueSign()]; !ok {
		c.failures[req.UniqueSign()] = req
		c.scheduler.Push(req)
	}
	// TODO: failed twice or more, adds req to failure queue
}
