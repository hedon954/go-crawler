package collect

// Request defines a crawler request
type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte) ParseResult
}

// ParseResult defines the result after parsing crawled response
type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
