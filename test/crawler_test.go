package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedon954/go-crawler/collector/sqlstorage"
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

	// fetcher
	f = &fetcher.BrowserFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	// storage
	storage, err = sqlstorage.New(
		sqlstorage.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(l.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	)

	// limiter
	secondLimit  = rate.NewLimiter(limiter.Per(1, 2*time.Second), 1)
	minuteLimit  = rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)
	multiLimiter = limiter.NewMultiLimiter(secondLimit, minuteLimit)
)

func TestCrawler_Run_WithStorage(t *testing.T) {
	l.Info("log init end")
	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}
	seeds := make([]*fetcher.Task, 0, 1000)
	seeds = append(seeds, &fetcher.Task{
		Property: fetcher.Property{
			Name:   douban.TaskNameDoubanBook,
			Cookie: `bid=j6xivD5rotM; _pk_ses.100001.3ac3=*; ap_v=0,6.0; __utma=30149280.1159442494.1678959267.1678959267.1678959267.1; __utmc=30149280; __utmz=30149280.1678959267.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmt_douban=1; __utma=81379588.1008421205.1678959267.1678959267.1678959267.1; __utmc=81379588; __utmz=81379588.1678959267.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmt=1; __ya…`,
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
