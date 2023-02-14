package fetcher

import (
	"regexp"
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
	ParseFunc func(*Context) (ParseResult, error)
}

type Context struct {
	Body []byte
	Req  *Request
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
