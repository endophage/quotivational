package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
)

// NewAuthHandler returns a server handler for authentication
func NewAuthHandler(goodTokens []string) http.Handler {
	goodTokensByID := make(map[string]bool)
	for _, token := range goodTokens {
		goodTokensByID[token] = true
	}

	m := mux.NewRouter()
	m.Methods("GET").Path("/token/{token:.+}").Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			if _, ok := goodTokensByID[vars["token"]]; !ok {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}))

	return m
}

func main() {
	flag.Parse()
	http.ListenAndServe(":8080", NewAuthHandler(flag.Args()))
}
