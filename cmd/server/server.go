package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Quote encapsulates a quote
type Quote struct {
	ID      int64     `xorm:"id"`
	Text    string    `xorm:"text"`
	Author  string    `xorm:"author"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

// QuoteServer sets up the quote server
type QuoteServer struct {
	db       *xorm.Engine
	redis    redis.Conn
	authaddr string
}

// NewQuoteServer is a constructor for QuoteServer
func NewQuoteServer(db *xorm.Engine, redisConn redis.Conn, authaddr string) *QuoteServer {
	return &QuoteServer{db: db, redis: redisConn,
		authaddr: strings.TrimSuffix(authaddr, "/")}
}

// GetQuoteHandler is the handler that returns the quotes
func (s *QuoteServer) GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("x-auth-token")
	authed, err := s.Authenticate(key)
	if err != nil {
		fmt.Println("error authenticating: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authed {
		fmt.Println("unauthorized: ", key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	quoteID, err := strconv.ParseInt(vars["quoteID"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	s.returnQuoteByID(w, quoteID)
}

func (s *QuoteServer) returnQuoteByID(w http.ResponseWriter, quoteID int64) {
	quote := &Quote{}
	has, err := s.db.Id(quoteID).Get(quote)

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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetRandomQuoteHandler is the handler that looks up what quotes have been
// seen so far and returns one that hasn't been seen lately
func (s *QuoteServer) GetRandomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("x-auth-token")
	authed, err := s.Authenticate(key)
	if err != nil {
		fmt.Println("error authenticating: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authed {
		fmt.Println("unauthorized: ", key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// watch/get/multi/exec - loop until success
	// see example at https://github.com/garyburd/redigo/blob/master/redis/zpop_example_test.go
	quoteID, err := redis.Int64(s.redis.Do("SPOP", key))
	if err == redis.ErrNil {
		// no more elements left in the set, so start over from the
		// beginning - we want to return the first quote but also
		// repopulate the cache with all the quote IDs
		defer func() {
			// Return connection to normal state on error.
			if err != nil {
				s.redis.Do("DISCARD")
			}
		}()

		var queued interface{}
		for {
			if err = s.redis.Send("WATCH", key); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			q := []Quote{}
			err := s.db.Desc("ID").Find(&q)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if len(q) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			maxID := q[0].ID
			minID := q[len(q)-1].ID

			s.redis.Send("MULTI")
			for i := minID + 1; i <= maxID; i++ {
				s.redis.Send("SADD", key, i)
			}

			queued, err = s.redis.Do("EXEC")
			if err != nil && err != redis.ErrNil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if queued != nil {
				quoteID = minID
				break
			}
		}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.returnQuoteByID(w, quoteID)
}

// Authenticate returns true if the request is authenticated, false else
func (s *QuoteServer) Authenticate(authToken string) (bool, error) {
	if authToken == "" {
		return false, nil
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/token/%s",
		s.authaddr, authToken), nil)
	if err != nil {
		return false, err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// ServerHandlers returns HTTP handlers for the server
func (s *QuoteServer) ServerHandlers() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/quotes/{quoteID:[0-9]+}").Handler(
		http.HandlerFunc(s.GetQuoteHandler))
	r.Methods("GET").Path("/randomquote").Handler(
		http.HandlerFunc(s.GetRandomQuoteHandler))
	return r
}

func setupSQL(dbtype, dbsource string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(dbtype, dbsource)
	if err != nil {
		return nil, err
	}
	err = engine.CreateTables(&Quote{})
	if err != nil {
		engine.Close()
		return nil, err
	}
	return engine, nil
}

func main() {
	var mysqldb = flag.String("db", "", "The DB source")
	var redisAddr = flag.String("redis", "", "Where Redis is")
	var authserver = flag.String("auth", "", "Where the auth server is")

	flag.Parse()

	var engine *xorm.Engine
	var redisConn redis.Conn
	var err error

	for {
		if engine == nil {
			engine, err = setupSQL("mysql", *mysqldb)
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
		}

		redisConn, err = redis.Dial("tcp", *redisAddr)
		if err == nil {
			break
		}
		fmt.Println(err.Error())
		time.Sleep(5 * time.Second)
	}
	defer engine.Close()
	defer redisConn.Close()

	q := NewQuoteServer(engine, redisConn, *authserver)
	fmt.Println("Starting server")
	http.ListenAndServe(":8080", q.ServerHandlers())
}
