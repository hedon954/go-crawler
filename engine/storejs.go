package engine

import (
	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/parser/douban"
	"github.com/robertkrimen/otto"
)

func init() {
	Store.AddJSTask(douban.DoubanJsTask)
}

// AddJSReqs 添加 js 动态爬取规则
func AddJSReqs(jreqs []map[string]interface{}) []*fetcher.Request {
	reqs := make([]*fetcher.Request, 0)
	for _, jreq := range jreqs {
		req := &fetcher.Request{}
		url, ok := jreq["Url"].(string)
		if !ok {
			continue
		}
		req.Url = url
		req.RuleName, _ = jreq["RuleName"].(string)
		req.Method, _ = jreq["Method"].(string)
		req.Priority, _ = jreq["Priority"].(int64)
		reqs = append(reqs, req)
	}
	return reqs
}

// AddJSReq 添加 js 动态爬取规则
func AddJSReq(jreq map[string]interface{}) []*fetcher.Request {
	reqs := make([]*fetcher.Request, 0)
	req := &fetcher.Request{}
	url, ok := jreq["Url"].(string)
	if !ok {
		return nil
	}
	req.Url = url
	req.RuleName, _ = jreq["RuleName"].(string)
	req.Method, _ = jreq["Method"].(string)
	req.Priority, _ = jreq["Priority"].(int64)
	reqs = append(reqs, req)
	return reqs
}

// AddJSTask 添加 js 动态爬取任务
func (cs *CrawlerStore) AddJSTask(m *fetcher.TaskModel) {
	task := &fetcher.Task{
		Property: m.Property,
	}
	task.Rule.Root = func() ([]*fetcher.Request, error) {
		vm := otto.New()
		err := vm.Set("AddJsReq", AddJSReqs)
		if err != nil {
			return nil, err
		}
		v, err := vm.Eval(m.Root)
		if err != nil {
			return nil, err
		}
		e, err := v.Export()
		if err != nil {
			return nil, err
		}
		return e.([]*fetcher.Request), nil
	}

	for _, r := range m.Rules {
		parseFunc := func(parse string) func(ctx *fetcher.Context) (fetcher.ParseResult, error) {
			return func(ctx *fetcher.Context) (fetcher.ParseResult, error) {
				vm := otto.New()
				err := vm.Set("ctx", ctx)
				if err != nil {
					return fetcher.ParseResult{}, err
				}
				v, err := vm.Eval(parse)
				if err != nil {
					return fetcher.ParseResult{}, err
				}
				e, err := v.Export()
				if err != nil {
					return fetcher.ParseResult{}, err
				}
				if e == nil {
					return fetcher.ParseResult{}, err
				}
				return e.(fetcher.ParseResult), err
			}
		}(r.ParseFunc)
		if task.Rule.Trunk == nil {
			task.Rule.Trunk = make(map[string]*fetcher.Rule, 0)
		}
		task.Rule.Trunk[r.Name] = &fetcher.Rule{
			ParseFunc: parseFunc,
		}
	}

	cs.Hash[task.Name] = task
	cs.list = append(cs.list, task)
}
