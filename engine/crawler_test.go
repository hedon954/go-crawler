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

func TestCrawler_Run(t *testing.T) {
	plugin := logger.NewStdoutPlugin(zapcore.InfoLevel)
	l := logger.NewLogger(plugin)
	l.Info("log init end")

	var f fetcher.Fetcher = &fetcher.BrowserFetcher{
		Timeout: 3000 * time.Millisecond,
		Logger:  l,
	}

	// douban cookie
	var seeds = make([]*fetcher.Task, 0, 1000)
	for i := 0; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &fetcher.Task{
			Url:      str,
			WaitTime: 1 * time.Second,
			MaxDepth: 5,
			Fetcher:  f,
			Cookie:   "bid=-UXUw--yL5g; dbcl2=\"214281202:q0BBm9YC2Yg\"; __yadk_uid=jigAbrEOKiwgbAaLUt0G3yPsvehXcvrs; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1665849857.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.21428; ck=SAvm; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1665925405%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=*; __utma=30149280.2072705865.1665849857.1665849857.1665925407.2; __utmc=30149280; __utmt=1; __utmb=30149280.23.5.1665925419338; _pk_id.100001.8cb4=fc1581490bf2b70c.1665849856.2.1665925421.1665849856.",
			RootReq: &fetcher.Request{
				ParseFunc: douban.ParseCityList,
				Priority:  i,
			},
		})
	}

	s := NewCrawler(
		WithFetcher(f),
		WithLogger(l),
		WithWorkCount(5),
		WithSeeds(seeds),
		WithScheduler(NewSchedule()),
	)

	s.Run()

}
