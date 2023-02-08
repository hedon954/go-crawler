package main

import (
	"fmt"

	"github.com/hedon954/go-crawler/collect"
)

const (
	urlBaidu  = "http://www.baidu.com"
	urlDouban = "https://book.douban.com/subject/1007305/"
)

func main() {
	url := urlDouban
	bf := collect.BrowserFetcher{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}
	bs, err := bf.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
