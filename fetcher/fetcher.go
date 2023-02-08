// Package fetcher
// @description implements a crawler fetcher
package fetcher

// Fetcher defines the crawler engine behaviors
type Fetcher interface {

	// Fetch the html content according to url
	Get(url *Request) ([]byte, error)
}
