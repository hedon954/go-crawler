package douban

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedon954/go-crawler/fetcher"
	"github.com/hedon954/go-crawler/logger"
	"go.uber.org/zap"
)

func TestParseCityList(t *testing.T) {
	plugin, closer := logger.NewFilePlugin("test.log", zap.DebugLevel)
	defer closer.Close()
	l := logger.NewLogger(plugin)

	var workList []*fetcher.Request
	for i := 0; i < 100; i += 25 {
		str := fmt.Sprintf(originUrl, i)
		workList = append(workList, &fetcher.Request{
			Url:       str,
			ParseFunc: ParseCityList,
		})
	}

	fetcher := fetcher.BrowserFetcher{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			bs, err := fetcher.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				l.Error("read content failed", zap.Error(err))
				continue
			}
			res := item.ParseFunc(bs, item)
			for _, i := range res.Items {
				l.Info("resuult", zap.String("get url:", i.(string)))
			}
			workList = append(workList, res.Requests...)
		}
	}
}
