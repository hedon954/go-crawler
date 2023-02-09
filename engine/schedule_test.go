package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/logger"
	"github.com/hedon954/go-crawler/parser/douban"
	"go.uber.org/zap/zapcore"
)

const (
	doubanUrl = "https://www.douban.com/group/szsh/discussion?start=%d"
)

func TestScheduler_Run(t *testing.T) {

	plugin := logger.NewStderrPlugin(zapcore.DebugLevel)
	l := logger.NewLogger(plugin)

	var seeds []*fetcher.Task
	for i := 0; i <= 100; i += 25 {
		url := fmt.Sprintf(doubanUrl, i)
		seeds = append(seeds, &fetcher.Task{
			Url:      url,
			WaitTime: 3 * time.Second,
			Cookie:   "xxx",
			MaxDepth: 5,
			RootReq: &fetcher.Request{
				ParseFunc: douban.ParseCityList,
			},
		})
	}

	f := &fetcher.BrowserFetcher{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		Timeout:   3 * time.Second,
	}

	e := NewScheduler(
		WithWorkCount(5),
		WithFetcher(f),
		WithLogger(l),
		WithSeeds(seeds),
	)
	e.Run()
}
