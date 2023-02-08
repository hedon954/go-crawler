// Package collect
// @description implements a crawler engine
package collect

// Fetcher defines the crawler engine behaviors
type Fetcher interface {

	// Fetch the html content according to url
	Get(url *Request) ([]byte, error)
}
