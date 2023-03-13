package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedon954/go-crawler/collector/sqlstorage"
	"github.com/hedon954/go-crawler/engine"
	"github.com/hedon954/go-crawler/parser/douban"

	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/logger"
	"go.uber.org/zap/zapcore"
)

func TestCrawler_Run_WithStorage(t *testing.T) {
	plugin := logger.NewStdoutPlugin(zapcore.InfoLevel)
	l := logger.NewLogger(plugin)
	l.Info("log init end")

	var f fetcher.Fetcher = &fetcher.BrowserFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	storage, err := sqlstorage.New(
		sqlstorage.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(l.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	)
	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}

	seeds := make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Property: fetcher.Property{
			Name: douban.TaskNameDoubanBook,
		},
		Fetcher: f,
		Storage: storage,
	})

	s := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(l),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	s.Run()
}
