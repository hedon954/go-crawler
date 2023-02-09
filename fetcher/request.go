package fetcher

import (
	"errors"
	"time"
)

// Request defines a crawler request
type Request struct {
	Url     string
	Cookie  string
	Timeout time.Duration
	// Current crawling depth
	Depth int
	// The max crawling depth
	MaxDepth  int
	ParseFunc func([]byte, *Request) ParseResult
}

// ParseResult defines the result after parsing crawled response
type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

// Check 对 request 进行合法性检查
func (r Request) Check() error {
	if err := r.checkDepth(); err != nil {
		return err
	}
	return nil
}

// checkDepth 检查深度的合法性
func (r Request) checkDepth() error {
	if r.Depth > r.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}
