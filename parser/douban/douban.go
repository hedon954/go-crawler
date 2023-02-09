package douban

import (
	"regexp"

	"github.com/hedon954/go-crawler/fetcher"
)

const (
	originUrl = "https://www.douban.com/group/szsh/discussion?start=%d"

	cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
	contentRe  = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`
)

// ParseCityList parses douban city url list
func ParseCityList(contents []byte, req *fetcher.Request) fetcher.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := fetcher.ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		result.Requests = append(result.Requests, &fetcher.Request{
			Task:  req.Task,
			Url:   url,
			Depth: req.Depth + 1,
			ParseFunc: func(c []byte, r *fetcher.Request) fetcher.ParseResult {
				return getContent(c, url)
			},
		})
	}

	return result
}

// getContent gets the url in content
func getContent(content []byte, url string) fetcher.ParseResult {
	re := regexp.MustCompile(contentRe)
	if ok := re.Match(content); !ok {
		return fetcher.ParseResult{}
	}

	return fetcher.ParseResult{
		Items: []interface{}{url},
	}
}
