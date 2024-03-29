package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hedon954/go-crawler/extensions"
	"github.com/hedon954/go-crawler/proxy"
	"go.uber.org/zap"
	"golang.org/x/text/transform"
)

// BrowserFetcher is a fetcher which simulates browser
type BrowserFetcher struct {
	Timeout time.Duration
	Proxy   proxy.Func
	Logger  *zap.Logger
}

func (b BrowserFetcher) Get(r *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	// Set proxy
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		return nil, err
	}

	// Set the header of User-Agent to simulate browser
	req.Header.Set("User-Agent", extensions.GenerateRandomUA())
	// Set cookie to simulate login status
	if len(r.Task.Cookie) > 0 {
		req.Header.Set("Cookie", r.Task.Cookie)
	}
	// Add specify headers
	for k, v := range r.Task.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	b.Logger.Info("start to fetch: " + r.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error status code: %d, status: %s\n", resp.StatusCode, resp.Status)
	}
	bodyReader := bufio.NewReader(resp.Body)
	encodeMode := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
