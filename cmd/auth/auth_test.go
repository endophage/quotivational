package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRequest(t *testing.T, baseURL, token string, expected int) {
	req, err := http.NewRequest("GET", baseURL+"/token/"+token, nil)
	if err != nil {
		t.Fatalf("expected no error setting up a request: %s", err)
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("should not have gotten an error making a request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != expected {
		t.Fatalf("expected a %v response, got %v", expected, resp.StatusCode)
	}
}

func TestAuthHandler(t *testing.T) {
	goodTokens := []string{"goody", "two", "shoes"}
	badTokens := []string{"baddy", "three", "boots"}

	s := httptest.NewServer(NewAuthHandler(goodTokens))

	for _, tok := range goodTokens {
		makeRequest(t, s.URL, tok, http.StatusOK)
	}

	for _, tok := range badTokens {
		makeRequest(t, s.URL, tok, http.StatusUnauthorized)
	}
}
