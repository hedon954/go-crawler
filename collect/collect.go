// Package collect
// @description implements a crawler engine
package collect

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Fetcher defines the crawler engine behaviors
type Fetcher interface {

	// Fetch the html content according to url
	Get(url string) ([]byte, error)
}

// BaseFetch is a basic implementation of Fetcher
type BaseFetch struct{}

func (bf BaseFetch) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("got error status code: %d", resp.StatusCode)
		// Do not return right now, just continue to read the response
	}
	bodyReader := bufio.NewReader(resp.Body)

	// In the Go language, strings are encoded in UTF-8 by default.
	// So here we always convert html content to UTF-8 format
	encodeMode := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// DeterminEncoding returns the encoder of the html content
func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bs, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("peek body error: %v\n", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bs, "")
	return e
}
