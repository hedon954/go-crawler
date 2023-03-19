package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hedon954/go-crawler/collector/sqlstorage"

	"github.com/hedon954/go-crawler/collector/gorm"

	"github.com/hedon954/go-crawler/parser/tianyancha"

	"github.com/hedon954/go-crawler/engine"
	"github.com/hedon954/go-crawler/limiter"
	"github.com/hedon954/go-crawler/parser/douban"
	"golang.org/x/time/rate"

	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/logger"
	"go.uber.org/zap/zapcore"
)

var (
	// logger
	plugin = logger.NewStdoutPlugin(zapcore.DebugLevel)
	l      = logger.NewLogger(plugin)

	// limiter
	secondLimit  = rate.NewLimiter(limiter.Per(60, time.Second), 60)
	minuteLimit  = rate.NewLimiter(limiter.Per(3000, 1*time.Minute), 3000)
	multiLimiter = limiter.NewMultiLimiter(secondLimit, minuteLimit)
)

func TestCrawler_Run_TianYanCha2(t *testing.T) {

	// storage
	storage, err := gorm.New(
		gorm.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		gorm.WithLogger(l.Named("sqlDB")),
		gorm.WithBatchCount(100),
	)

	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}
	go func() {
		lastCount := 0
		sameContinue := 0
		for {
			time.Sleep(1 * time.Minute)
			count, _ := storage.Flush()
			if lastCount == count {
				sameContinue++
				if sameContinue >= 2 {
					os.Exit(1)
					return
				}
			} else {
				sameContinue = 0
				lastCount = count
			}
		}
	}()

	// fetcher
	f := &fetcher.RedirectFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	l.Info("log init end")

	seeds := make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Property: fetcher.Property{
			Name: tianyancha.TaskNameTianYanCha2,
		},
		Fetcher: f,
		Storage: storage,
		Limiter: *multiLimiter,
	})

	s := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(l),
		engine.WithWorkCount(20),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	s.Run()
}

func TestCrawler_Run_TianYanCha(t *testing.T) {

	// storage
	storage, err := gorm.New(
		gorm.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		gorm.WithLogger(l.Named("sqlDB")),
		gorm.WithBatchCount(100),
	)

	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}
	go func() {
		lastCount := 0
		sameContinue := 0
		for {
			time.Sleep(1 * time.Minute)
			count, _ := storage.Flush()
			if lastCount == count {
				sameContinue++
				if sameContinue >= 2 {
					os.Exit(1)
					return
				}
			} else {
				sameContinue = 0
				lastCount = count
			}
		}
	}()

	// fetcher
	f := &fetcher.RedirectFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	l.Info("log init end")

	seeds := make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Property: fetcher.Property{
			Name:   tianyancha.TaskNameTianYanCha,
			Cookie: "xxx",
			Headers: map[string]string{
				"X-AUTH-TOKEN": "xxx",
				"X-TYCID":      "xxx",
			},
		},
		Fetcher: f,
		Storage: storage,
		Limiter: *multiLimiter,
	})

	s := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(l),
		engine.WithWorkCount(20),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	s.Run()
}

func TestCrawler_Run_WithStorage(t *testing.T) {

	// storage
	storage, err := sqlstorage.New(
		sqlstorage.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(l.Named("sqlDB")),
		sqlstorage.WithBatchCount(10),
	)

	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}

	// fetcher
	f := &fetcher.BrowserFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	l.Info("log init end")

	seeds := make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Property: fetcher.Property{
			Name:   douban.TaskNameDoubanBook,
			Cookie: "xxx",
		},
		Fetcher: f,
		Storage: storage,
		Limiter: *multiLimiter,
	})

	s := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(l),
		engine.WithWorkCount(10),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	s.Run()
}
