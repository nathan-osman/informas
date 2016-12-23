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
	sessions    *sessions.CookieStore
	config      *db.Config
	templateDir string
}

// New creates a new server instance.
func New(addr, dataDir string) (*Server, error) {
	c, err := db.NewConfig(&db.Token{})
	if err != nil {
		return nil, err
	}
	var (
		h = []byte(c.GetString(configSecretKey))
		m = mux.NewRouter()
		s = &Server{
			server:      server.New(addr),
			sessions:    sessions.NewCookieStore(h),
			config:      c,
			templateDir: path.Join(dataDir, "templates"),
		}
	)
	s.server.Handler = m
	m.HandleFunc("/", s.view(accessRegistered, s.index))
	m.HandleFunc("/install", s.view(accessPublic, s.install))
	m.HandleFunc("/login", s.view(accessPublic, s.login))
	m.HandleFunc("/logout", s.view(accessRegistered, s.logout))
	m.PathPrefix("/static").Handler(
		http.FileServer(http.Dir(dataDir)),
	)
	return s, nil
}

// Start begins listening on the specified address.
func (s *Server) Start() error {
	return s.server.Start()
}

// Stop shuts down the server.
func (s *Server) Stop() {
	s.server.Stop()
}
