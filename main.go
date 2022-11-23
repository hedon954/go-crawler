package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fetch url error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error status code: %v", resp.StatusCode)
		return
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response body failed: %v", err)
		return
	}

	fmt.Println("body:", string(bs))
}
