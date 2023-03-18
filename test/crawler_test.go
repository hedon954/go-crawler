package test

import (
	"fmt"
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
	secondLimit  = rate.NewLimiter(limiter.Per(1, 2*time.Second), 1)
	minuteLimit  = rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)
	multiLimiter = limiter.NewMultiLimiter(secondLimit, minuteLimit)
)

func TestCrawler_Run_TianYanCha(t *testing.T) {

	// storage
	storage, err := gorm.New(
		gorm.WithSqlUrl("root:root@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		gorm.WithLogger(l.Named("sqlDB")),
		gorm.WithBatchCount(10),
	)

	if err != nil {
		l.Error(fmt.Sprintf("create sqlstorage failed: %v", err))
		return
	}

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
			Cookie: `HWWAFSESID=9edbfc72abc96c645b; HWWAFSESTIME=1679110175285; csrfToken=wDGWVZ58dxKLHKD0RJplBPEV; jsid=SEO-BAIDU-ALL-SY-000001; TYCID=1af08280c53d11ed81b4bdbf95e432f8; sajssdk_2015_cross_new_user=1; Hm_lvt_e92c8d65d92d534b0fc290df538b4758=1679110176; bdHomeCount=0; bannerFlag=true; ssuid=8176940507; _ga=GA1.2.292669649.1679110386; _gid=GA1.2.1479135818.1679110386; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22297428115%22%2C%22first_id%22%3A%22186f2c3fbd1926-076f09e7232fd44-1f525634-1296000-186f2c3fbd2b8e%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTg2ZjJjM2ZiZDE5MjYtMDc2ZjA5ZTcyMzJmZDQ0LTFmNTI1NjM0LTEyOTYwMDAtMTg2ZjJjM2ZiZDJiOGUiLCIkaWRlbnRpdHlfbG9naW5faWQiOiIyOTc0MjgxMTUifQ%3D%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%24identity_login_id%22%2C%22value%22%3A%22297428115%22%7D%2C%22%24device_id%22%3A%22186f2c3fbd1926-076f09e7232fd44-1f525634-1296000-186f2c3fbd2b8e%22%7D; tyc-user-info={%22state%22:%220%22%2C%22vipManager%22:%220%22%2C%22mobile%22:%2215623205156%22}; tyc-user-info-save-time=1679110440917; auth_token=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxNTYyMzIwNTE1NiIsImlhdCI6MTY3OTExMDQ0MCwiZXhwIjoxNjgxNzAyNDQwfQ.mNlC3Kj2Y5OKRozqh59wOjpE1jLuCKtRVtqcXdKeh_IdB7Ym9tfTjey-2x7gjxfZRwcZ2U76B74mp_ox2hT9hQ; tyc-user-phone=%255B%252215623205156%2522%255D; searchSessionId=1679116419.23758314; Hm_lpvt_e92c8d65d92d534b0fc290df538b4758=1679116424; cloud_token=535afaea9df647bd8f27a3356324317b; cloud_utm=4cc2926538634b1a896db92fa9034cda`,
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
