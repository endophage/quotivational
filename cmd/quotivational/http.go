package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	authHeader = "X-Auth-Token"
	randomPath = "/randomquote"
	topicPath  = "/quotes/%s"
)

// HTTPQuoter gets a quote from a quote server
type HTTPQuoter struct {
	url   url.URL
	token string
}

// NewHTTPQuoter returns a HTTPQuoter instance
func NewHTTPQuoter(baseURL string, authToken string) (*HTTPQuoter, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &HTTPQuoter{
		url:   *u,
		token: authToken,
	}, nil
}

// Quote gets a quote of a particular topic
func (h HTTPQuoter) Quote(topic string) (*Quote, error) {
	u, err := h.url.Parse(
		fmt.Sprintf(topicPath, topic),
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header[authHeader] = []string{h.token}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	q := &Quote{}
	err = json.Unmarshal(b, q)
	if err != nil {
		return nil, err
	}
	return q, nil
}
