package main

import (
	"fmt"

	"github.com/hedon954/go-crawler/collect"
)

func main() {
	url := "http://www.baidu.com"
	bf := collect.BaseFetch{}
	bs, err := bf.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
