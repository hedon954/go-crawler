package collect

import (
	"bufio"
	"io/ioutil"
	"net/http"

	"golang.org/x/text/transform"
)

// BrowserFetcher is a fetcher which simulates browser
type BrowserFetcher struct {
	UserAgent string
}

func (b BrowserFetcher) Get(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the header of User-Agent to simulate browser
	req.Header.Set("User-Agent", b.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyReader := bufio.NewReader(resp.Body)
	encodeMode := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
