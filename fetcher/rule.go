package fetcher

import (
	"regexp"
	"time"

	"github.com/hedon954/go-crawler/collector"
)

type RuleTree struct {

	// the entry of crawling rules
	Root func() ([]*Request, error)

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

// ParseJSReg 用于 JS 代码中解析正则表达式，获取请求任务列表
func (c *Context) ParseJSReg(ruleName string, reg string) ParseResult {
	re := regexp.MustCompile(reg)
	matches := re.FindAllSubmatch(c.Body, -1)

	res := ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		res.Requests = append(res.Requests, &Request{
			Method:   "GET",
			Task:     c.Req.Task,
			Url:      url,
			Depth:    c.Req.Depth + 1,
			RuleName: ruleName,
		})
	}
	return res
}

// OutputJS 用于 JS 代码中解析正则表达式，获取爬取结果
func (c *Context) OutputJS(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	ok := re.Match(c.Body)
	if !ok {
		return ParseResult{Items: []interface{}{}}
	}
	return ParseResult{
		Items: []interface{}{c.Req.Url},
	}
}
