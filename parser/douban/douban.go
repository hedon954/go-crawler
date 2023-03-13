package douban

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hedon954/go-crawler/fetcher"
)

const (
	originUrl = "https://www.douban.com/group/szsh/discussion?start=%d"

	cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
	contentRe  = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

	ruleNameParseUrl       = "parse_website_url"
	ruleNameResolveSunRoom = "resolve_sun_room"

	TaskNameFindSunRoom = "find_douban_sun_room"
)

var DoubanTask = &fetcher.Task{
	Property: fetcher.Property{
		Name:     TaskNameFindSunRoom,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "xxx",
	},
	Rule: fetcher.RuleTree{
		Root: func() ([]*fetcher.Request, error) {
			var roots []*fetcher.Request
			for i := 0; i < 100; i += 25 {
				url := fmt.Sprintf(originUrl, i)
				roots = append(roots, &fetcher.Request{
					Priority: 1,
					Url:      url,
					Method:   "GET",
					RuleName: ruleNameParseUrl,
				})
			}
			return roots, nil
		},
		Trunk: map[string]*fetcher.Rule{
			ruleNameParseUrl:       {nil, ParseURL},
			ruleNameResolveSunRoom: {nil, GetSunRoom},
		},
	},
}

// ParseURL parses the target url
func ParseURL(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)

	res := fetcher.ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		res.Requests = append(res.Requests, &fetcher.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      url,
			Depth:    ctx.Req.Depth + 1,
			RuleName: ruleNameResolveSunRoom,
		})
	}
	return res, nil
}

// GetSunRoom resolves the sun room infos in the crawled website
func GetSunRoom(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	re := regexp.MustCompile(contentRe)
	ok := re.Match(ctx.Body)
	if !ok {
		return fetcher.ParseResult{
			Items: []interface{}{},
		}, nil
	}
	return fetcher.ParseResult{
		Items: []interface{}{ctx.Req.Url},
	}, nil
}
