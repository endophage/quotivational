package main

var (
	allTopics = []string{"life", "computers", "science", "drinking"}
)

type Quoter interface {
	Quote(topic string) (Quote, error)
}

type Quote struct {
	Quote  string `json:"Text"`
	Author string `json:"Author"`
}
