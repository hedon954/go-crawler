package proxy

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProxy(t *testing.T) {
	var index uint32 = 1
	var urlCounts uint32 = 10

	r := roundRobinSwitcher{}
	r.index = index
	r.proxyURLs = make([]*url.URL, urlCounts)

	for i := 0; i < int(urlCounts); i++ {
		r.proxyURLs[i] = &url.URL{}
		r.proxyURLs[i].Host = strconv.Itoa(i)
	}

	p, err := r.GetProxy(nil)
	if err != nil && strings.Contains(err.Error(), "empty proxy urls") {
		t.Skip()
	}
	assert.Nil(t, err)

	e := r.proxyURLs[index&urlCounts]
	if !reflect.DeepEqual(p, e) {
		t.Fail()
	}
}
