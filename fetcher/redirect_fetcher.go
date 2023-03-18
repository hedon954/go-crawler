package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hedon954/go-crawler/extensions"
	"github.com/hedon954/go-crawler/proxy"
	"go.uber.org/zap"
	"golang.org/x/text/transform"
)

// RedirectFetcher is a fetcher that deals with redirected links
type RedirectFetcher struct {
	Timeout time.Duration
	Proxy   proxy.Func
	Logger  *zap.Logger
}

func (b RedirectFetcher) Get(r *Request) ([]byte, error) {

	// Sometimes the home page of some websites will use delayed loading,
	// and sending a GET request directly may not bring up all the information.
	// At this point, you can save them as local files,
	// read the local files, and then do dynamic parsing later.
	if bs, err, ok := b.handleLocalFile(r); ok {
		return bs, err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
	req.Header.Set("allow_redirects", "true")

	// Set the header of User-Agent to simulate browser
	req.Header.Set("User-Agent", extensions.GenerateRandomUA())
	// Set cookie to simulate login status
	if len(r.Task.Cookie) > 0 {
		req.Header.Set("Cookie", r.Task.Cookie)
	}

	resp, err := client.Do(req)
	b.Logger.Info("start to fetch: " + r.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// redirect, get the real url
	if resp.StatusCode == http.StatusTemporaryRedirect {
		realPath := resp.Header.Get("Location")
		if realPath != "" {
			urlInfo, _ := url.Parse(r.Url)
			host := urlInfo.Host
			newUrl := "https://" + host + realPath
			return b.handleRedirectUrl(newUrl, r.Task.Cookie)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error status code: %d, status: %s\n", resp.StatusCode, resp.Status)
	}
	bodyReader := bufio.NewReader(resp.Body)
	encodeMode := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// handleRedirectUrl sends http request get the redirect url body
func (b RedirectFetcher) handleRedirectUrl(newUrl, cookie string) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	// Set proxy
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return nil, err
	}

	// Set the header of User-Agent to simulate browser
	req.Header.Set("User-Agent", extensions.GenerateRandomUA())
	// Set cookie to simulate login status
	if len(cookie) > 0 {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := client.Do(req)
	b.Logger.Info("start to fetch the redirect url: " + newUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error status code: %d, status: %s\n", resp.StatusCode, resp.Status)
	}
	bodyReader := bufio.NewReader(resp.Body)
	encodeMode := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, encodeMode.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func (b RedirectFetcher) handleLocalFile(r *Request) ([]byte, error, bool) {
	if r.Url == "https://www.tianyancha.com/" {
		f, err := os.Open("../parser/tianyancha/tianyancha.html")
		if err != nil {
			return nil, err, true
		}
		bs, err := ioutil.ReadAll(f)
		return bs, err, true
	}

	return nil, nil, false
}
