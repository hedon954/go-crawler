package fetcher

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// BaseFetcher is a basic implementation of Fetcher
type BaseFetcher struct{}

func (bf BaseFetcher) Get(r *Request) ([]byte, error) {
	resp, err := http.Get(r.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("got error status code: %d, status: %s\n", resp.StatusCode, resp.Status)
		// Do not return right now, just continue to read the response
	}
	bodyReader := bufio.NewReader(resp.Body)

	// In the Go language, strings are encoded in UTF-8 by default.
	// So here we always convert html content to UTF-8 format
	encodeMode := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// DetermineEncoding returns the encoder of the html content
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bs, err := r.Peek(1024)
	if err != nil && err != io.EOF {
		fmt.Printf("peek body error: %v\n", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bs, "")
	return e
}
