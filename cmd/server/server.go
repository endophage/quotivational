package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Quote encapsulates a quote
type Quote struct {
	ID      int64     `xorm:"pk"`
	Text    string    `xorm:"text not null"`
	Author  string    `xorm:"varchar(255) not null"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

// QuoteServer sets up the quote server
type QuoteServer struct {
	db *xorm.Engine
}

// NewQuoteServer is a constructor for QuoteServer
func NewQuoteServer(db *xorm.Engine) *QuoteServer {
	return &QuoteServer{db: db}
}

// GetQuoteHandler is the handler that returns the quotes
func (s *QuoteServer) GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quoteID, err := strconv.ParseInt(vars["quoteID"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	quote := &Quote{ID: quoteID}
	has, err := s.db.Get(quote)

	if has && err == nil {
		var result []byte
		result, err = json.Marshal(quote)
		if err == nil {
			w.Write(result)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ServerHandlers returns HTTP handlers for the server
func (s *QuoteServer) ServerHandlers() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/quotes/{quoteID:[0-9]+}").Handler(
		http.HandlerFunc(s.GetQuoteHandler))
	return r
}

func fatalf(formatString string, args ...interface{}) {
	fmt.Printf(formatString, args...)
	os.Exit(1)
}

func setupSQlite(dbdir string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", filepath.Join(dbdir, "db"))
	if err != nil {
		return nil, err
	}
	err = engine.CreateTables(&Quote{})
	if err != nil {
		engine.Close()
		os.Remove(filepath.Join(dbdir, "db"))
		return nil, err
	}
	return engine, nil
}

// func main() {
// 	engine, _, err := setupSQlite()
// 	if err != nil {
// 		fatalf(err.Error())
// 	}
// 	q := QuoteServer{db: engine}
// }
