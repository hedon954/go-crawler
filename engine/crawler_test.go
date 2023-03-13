package engine

import (
	"testing"
	"time"

	"github.com/hedon954/go-crawler/parser/douban"

	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/logger"
	"go.uber.org/zap/zapcore"
)

func TestCrawler_Run(t *testing.T) {
	plugin := logger.NewStdoutPlugin(zapcore.InfoLevel)
	l := logger.NewLogger(plugin)
	l.Info("log init end")

	var f fetcher.Fetcher = &fetcher.BrowserFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	var seeds = make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Name:    douban.TaskNameFindSunRoom,
		Fetcher: f,
	})

	s := NewCrawler(
		WithFetcher(f),
		WithLogger(l),
		WithWorkCount(5),
		WithSeeds(seeds),
		WithScheduler(NewSchedule()),
	)

	s.Run()

}
