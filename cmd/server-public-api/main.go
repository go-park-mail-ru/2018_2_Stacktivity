package main

import (
	"2018_2_Stacktivity/cmd/server-public-api/session"
	"2018_2_Stacktivity/cmd/server-public-api/storage"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"path"
	"strings"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	sm    session.SessionManagerI
	users storage.UserStorageI
	log   *log.Logger
}

func NewServer(logger *log.Logger, db *sql.DB) *Server {
	return &Server{
		sm:    session.NewSessionManager(),
		log:   logger,
		users: storage.NewUserStorage(db),
	}
}

func main() {
	flag.Parse()
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetOutput(os.Stdout) // TODO write log in file

	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	srv := NewServer(logger, db)
	err = srv.users.Prepare()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Infof("starting server listening on %s", config.Port)
	err = http.ListenAndServe(config.Port, srv)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.log.Infoln(r.Method + " " + r.URL.Path)
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "user":
		srv.RouteUser(w, r)
	case "session":
		srv.RouteSession(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) RouteUser(w http.ResponseWriter, r *http.Request) {
	var head string
	head, _ = ShiftPath(r.URL.Path)
	if head == "" {
		switch r.Method {
		case http.MethodGet:
			srv.getUsers(w, r)
		case http.MethodPost:
			srv.createUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
	switch r.Method {
	case http.MethodGet:
		srv.getUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) RouteSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		srv.getSession(w, r)
	case http.MethodPost:
		srv.createSession(w, r)
	case http.MethodDelete:
		srv.deleteSession(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
