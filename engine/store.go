package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/parser/douban"
)

func init() {
	Store.Add(douban.DoubanTask)
}

// Store is the global CrawlerStore instance
var Store = &CrawlerStore{
	list: []*fetcher.Task{},
	hash: map[string]*fetcher.Task{},
}

// CrawlerStore scores the crawler tasks
type CrawlerStore struct {
	list []*fetcher.Task
	hash map[string]*fetcher.Task
}

// Add adds a task to the global crawler instance
func (cs *CrawlerStore) Add(task *fetcher.Task) {
	cs.hash[task.Name] = task
	cs.list = append(cs.list, task)
}
