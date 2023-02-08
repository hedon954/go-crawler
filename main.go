package main

import (
	"fmt"
	"time"

	"github.com/hedon954/go-crawler/collect"
	"github.com/hedon954/go-crawler/proxy"
)

const (
	urlBaidu  = "http://www.baidu.com"
	urlDouban = "https://book.douban.com/subject/1007305/"
)

func main() {
	url := urlDouban

	// Get proxy
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	proxy, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		panic(err)
	}

	// Creates a fetcher
	bf := collect.BrowserFetcher{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		Timeout:   time.Second * 10,
		Proxy:     proxy,
	}

	// Crawl
	bs, err := bf.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
