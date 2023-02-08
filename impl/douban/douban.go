package douban

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hedon954/go-crawler/collect"
	"github.com/hedon954/go-crawler/logger"
	"go.uber.org/zap"
)

const (
	originUrl = "https://www.douban.com/group/szsh/discussion?start=%d"

	cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
	contentRe  = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`
)

func Crawling() {

	plugin, closer := logger.NewFilePlugin("test.log", zap.DebugLevel)
	defer closer.Close()
	l := logger.NewLogger(plugin)

	workList := GetWorkList()
	fetcher := collect.BrowserFetcher{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		Timeout:   3 * time.Second,
	}
	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			bs, err := fetcher.Get(item.Url)
			time.Sleep(1 * time.Second)
			if err != nil {
				l.Error("read content failed", zap.Error(err))
				continue
			}
			res := item.ParseFunc(bs)
			for _, i := range res.Items {
				l.Info("resuult", zap.String("get url:", i.(string)))
			}
			workList = append(workList, res.Requests...)
		}
	}
}

func GetWorkList() []*collect.Request {
	var worklist []*collect.Request
	for i := 0; i < 100; i += 25 {
		str := fmt.Sprintf(originUrl, i)
		worklist = append(worklist, &collect.Request{
			Url:       str,
			ParseFunc: ParseCityList,
		})
	}
	return worklist
}

func ParseCityList(contents []byte) collect.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := collect.ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Url: url,
			ParseFunc: func(c []byte) collect.ParseResult {
				return GetContent(c, url)
			},
		})
	}

	return result
}

func GetContent(content []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(contentRe)
	if ok := re.Match(content); !ok {
		return collect.ParseResult{}
	}

	return collect.ParseResult{
		Items: []interface{}{url},
	}
}
