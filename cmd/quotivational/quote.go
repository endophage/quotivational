package main

var (
	allTopics = []string{"life", "computers", "science", "drinking"}
)

// Quoter is an interface for an object which returns a quote
type Quoter interface {
	Quote(topic string) (Quote, error)
}

// Quote is a structure representing a quote and its author
type Quote struct {
	Quote  string `json:"Text"`
	Author string `json:"Author"`
}
