package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/parser/douban"
	"github.com/hedon954/go-crawler/parser/tianyancha"
)

func init() {
	Store.Add(douban.DoubanTask)
	Store.Add(douban.DoubanBookTask)
	Store.Add(tianyancha.TianYanChaTask)
}

// Store is the global CrawlerStore instance
var Store = &CrawlerStore{
	list: []*fetcher.Task{},
	Hash: map[string]*fetcher.Task{},
}

// CrawlerStore scores the crawler tasks
type CrawlerStore struct {
	list []*fetcher.Task
	Hash map[string]*fetcher.Task
}

// Add adds a task to the global crawler instance
func (cs *CrawlerStore) Add(task *fetcher.Task) {
	cs.Hash[task.Name] = task
	cs.list = append(cs.list, task)
}

// GetFields returns fields by taskName and ruleName
func GetFields(taskName, ruleName string) []string {
	return Store.Hash[taskName].Rule.Trunk[ruleName].ItemFields
}
