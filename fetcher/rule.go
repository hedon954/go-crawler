package fetcher

import (
	"github.com/hedon954/go-crawler/collector"
	"time"
)

type RuleTree struct {

	// the entry of crawling rules
	Root func() []*Request

	// the hashmap of rules
	// key: rule's name
	// value: the specific rule
	Trunk map[string]*Rule
}

// Rule represents the rule corresponding to the request
type Rule struct {
	ItemFields []string
	ParseFunc  func(*Context) (ParseResult, error)
}

// Context is the crawling context
type Context struct {
	Body []byte
	Req  *Request
}

func (c *Context) Output(data interface{}) *collector.OutputData {
	res := &collector.OutputData{}
	res.Data = make(map[string]interface{})
	res.TaskName = c.Req.Task.Name
	res.RuleName = c.Req.RuleName
	res.Url = c.Req.Url
	res.Time = time.Now().Format("2006-01-02 15:04:05")
	res.Data = data
	return res
}
