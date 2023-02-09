package fetcher

import (
	"errors"
)

// Request represents a single crawler request
type Request struct {
	Task      *Task
	Url       string
	Depth     int
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
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}
