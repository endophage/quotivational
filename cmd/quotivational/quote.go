package main

type Quoter interface {
	Quote(topic string) (Quote, error)
}

type Quote struct {
	Quote  string
	Author string
}
