package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetQuoteByIDNoSuchQuote(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "quotivational")
	if err != nil {
		t.Fatalf("unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempDir)
	engine, err := setupSQlite(tempDir)
	if err != nil {
		t.Fatalf("expected no error setting up SQLite: %s", err)
	}
	defer engine.Close()

	q := QuoteServer{db: engine}
	ts := httptest.NewServer(q.ServerHandlers())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/quotes/1", nil)
	if err != nil {
		t.Fatalf("expected no error setting up a request: %s", err)
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("should not have gotten an error making a request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected a not found response, got %v", resp.StatusCode)
	}
}

func TestGetQuoteByIDExistingQuote(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "quotivational")
	if err != nil {
		t.Fatalf("unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempDir)
	engine, err := setupSQlite(tempDir)
	if err != nil {
		t.Fatalf("expected no error setting up SQLite: %s", err)
	}
	defer engine.Close()

	expected := &Quote{
		Text:   "this is a quote",
		Author: "iman author",
	}
	_, err = engine.Insert(expected)
	if err != nil {
		t.Fatalf("expected no error inserting into SQLite: %s", err)
	}

	q := QuoteServer{db: engine}
	ts := httptest.NewServer(q.ServerHandlers())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/quotes/0", nil)
	if err != nil {
		t.Fatalf("expected no error setting up a request: %s", err)
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("should not have gotten an error making a request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected a 200 response, got %v", resp.StatusCode)
	}

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %s", err)
	}

	readQuote := &Quote{}
	if err := json.Unmarshal(respJSON, readQuote); err != nil {
		t.Fatalf("could not parse response: %s", err)
	}

	if expected.Text != readQuote.Text && expected.Author != readQuote.Author {
		t.Fatalf("%v is not what was expected", readQuote)
	}
}

func TestGetQuoteByIDServerError(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "quotivational")
	if err != nil {
		t.Fatalf("unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempDir)
	engine, err := setupSQlite(tempDir)
	if err != nil {
		t.Fatalf("expected no error setting up SQLite: %s", err)
	}

	expected := &Quote{
		Text:   "this is a quote",
		Author: "iman author",
	}
	_, err = engine.Insert(expected)
	if err != nil {
		t.Fatalf("expected no error inserting into SQLite: %s", err)
	}

	q := QuoteServer{db: engine}
	ts := httptest.NewServer(q.ServerHandlers())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/quotes/0", nil)
	if err != nil {
		t.Fatalf("expected no error setting up a request: %s", err)
	}

	// kill the engine
	engine.Close()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("should not have gotten an error making a request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected a 500 response, got %v", resp.StatusCode)
	}
}

func TestSettingUpSQLTwiceIsFine(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "quotivational")
	if err != nil {
		t.Fatalf("unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	engine, err := setupSQlite(tempDir)
	if err != nil {
		t.Fatalf("expected no error setting up SQLite: %s", err)
	}

	first := &Quote{
		Text:   "this is a quote",
		Author: "iman author",
	}
	_, err = engine.Insert(first)
	if err != nil {
		t.Fatalf("expected no error inserting into SQLite: %s", err)
	}

	engine.Close()

	engine, err = setupSQlite(tempDir)
	if err != nil {
		t.Fatalf("expected no error setting up SQLite a second time: %s", err)
	}
	defer engine.Close()

	gotten := &Quote{}
	_, err = engine.Get(gotten)
	if err != nil {
		t.Fatalf("expected no error getting a row: %s", err)
	}

	if first.Text != gotten.Text && first.Author != gotten.Author {
		t.Fatalf("%v is not what was expected", gotten)
	}
}
