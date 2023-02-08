package main

import (
	"fmt"

	"github.com/hedon954/go-crawler/collect"
)

func main() {
	url := "https://book.douban.com/subject/1007305/"
	bf := collect.BaseFetcher{}
	bs, err := bf.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
