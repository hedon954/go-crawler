// Package proxy
// @description implements http proxy
package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

// Func returns proxy url for current http request
type Func func(r *http.Request) (*url.URL, error)

// RoundRobinProxySwitcher returns a proxy function
func RoundRobinProxySwitcher(proxyURLs ...string) (Func, error) {
	if len(proxyURLs) < 1 {
		return nil, errors.New("proxy URL list is empty")
	}
	urls := make([]*url.URL, len(proxyURLs))
	for i, u := range proxyURLs {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
	}
	return (&roundRobinSwitcher{
		proxyURLs: urls,
		index:     0,
	}).GetProxy, nil
}

// roundRobinSwitcher is a round-robin url switcher
type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]
	return u, nil
}
