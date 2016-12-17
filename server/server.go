package server

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hectane/go-asyncserver"
	"github.com/nathan-osman/informas/db"
)

const sessionName = "session"

// Server provides the web interface for the application.
type Server struct {
	server      *server.AsyncServer
	mux         *mux.Router
	sessions    *sessions.CookieStore
	config      *db.Config
	templateDir string
}

// New creates a new server instance.
func New(addr, dataDir string) (*Server, error) {
	c, err := db.NewConfig()
	if err != nil {
		return nil, err
	}
	s := &Server{
		server:      server.New(addr),
		mux:         mux.NewRouter(),
		sessions:    sessions.NewCookieStore([]byte(c.Get(configSecretKey))),
		config:      c,
		templateDir: path.Join(dataDir, "templates"),
	}
	s.server.Handler = s
	s.mux.HandleFunc("/", s.index)
	s.mux.HandleFunc("/login", s.login)
	s.mux.PathPrefix("/static").Handler(
		http.FileServer(http.Dir(dataDir)),
	)
	return s, nil
}

// ServeHTTP performs some preliminary processing before handing the request
// off to the router.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}

// Start begins listening on the specified address.
func (s *Server) Start() error {
	return s.server.Start()
}

// Stop shuts down the server.
func (s *Server) Stop() {
	s.server.Stop()
}
