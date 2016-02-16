package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
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

	q := NewQuoteServer(engine, nil)
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

	q := NewQuoteServer(engine, nil)
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

	if readQuote.ID != 1 {
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

	q := NewQuoteServer(engine, nil)
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

func TestGetRandomQuoteNoQuotes(t *testing.T) {
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

	key := ""

	c := redigomock.NewConn()
	c.Command("SADD", key, 0).Expect("ok")
	c.Command("WATCH", key).Expect("ok")
	c.Command("MULTI").Expect("ok")
	c.Command("EXEC").Expect("ok")
	c.Command("DISCARD").ExpectError(fmt.Errorf("this should never be called"))
	c.Command("SPOP", key).ExpectError(redis.ErrNil)

	q := NewQuoteServer(engine, c)
	ts := httptest.NewServer(q.ServerHandlers())
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/randomquote", nil)
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

func TestGetRandomQuoteExistingQuote(t *testing.T) {
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

	expected := Quote{
		Text:   "this is a quote",
		Author: "iman author",
	}
	for i := int64(0); i < 2; i++ {
		e := expected
		_, err = engine.Insert(&e)
		if err != nil {
			t.Fatalf("expected no error inserting into SQLite: %s", err)
		}
	}

	key := ""

	c := redigomock.NewConn()
	c.Command("SADD", key, 1).Expect("ok")
	c.Command("WATCH").Expect("ok")
	c.Command("MULTI").Expect("ok")
	c.Command("EXEC").ExpectError(redis.ErrNil).Expect("ok").Expect("ok")
	c.Command("DISCARD").Expect("ok")
	c.Command("SPOP", key).ExpectError(redis.ErrNil).Expect(int64(2)).ExpectError(redis.ErrNil)

	q := NewQuoteServer(engine, c)
	ts := httptest.NewServer(q.ServerHandlers())
	defer ts.Close()

	for i := int64(0); i < 3; i++ {
		fmt.Println("on i", i)
		req, err := http.NewRequest("GET", ts.URL+"/randomquote", nil)
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

		if readQuote.ID != i%2+1 {
			t.Fatalf("%v is not what was expected", readQuote)
		}
	}
}
