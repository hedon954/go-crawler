package fetcher

import (
	"time"
)

// Request defines a crawler request
type Request struct {
	Url       string
	Cookie    string
	Timeout   time.Duration
	ParseFunc func([]byte, *Request) ParseResult
}

// ParseResult defines the result after parsing crawled response
type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
