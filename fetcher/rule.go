package fetcher

type RuleTree struct {

	// the entry of crawling rules
	Root func() []*Request

	// the hashmap of rules
	// key: rule's name
	// value: the specific rule
	Trunk map[string]*Rule
}

// Rule represents the rule corresponding to the request
type Rule struct {
	ParseFunc func(*Context) ParseResult
}

type Context struct {
	Body []byte
	Req  *Request
}
