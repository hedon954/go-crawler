package fetcher

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// Request represents a single crawler request
type Request struct {
	unique   string
	Task     *Task
	Url      string
	Method   string
	Depth    int64
	Priority int64
	RuleName string
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

// UniqueSign builds the unique sign for each request
func (r Request) UniqueSign() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
