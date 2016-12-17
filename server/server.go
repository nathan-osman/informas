package server

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/hectane/go-asyncserver"
)

// Server provides the web interface for the application.
type Server struct {
	server      *server.AsyncServer
	mux         *mux.Router
	templateDir string
}

// New creates a new server instance.
func New(addr, dataDir string) (*Server, error) {
	s := &Server{
		server:      server.New(addr),
		mux:         mux.NewRouter(),
		templateDir: path.Join(dataDir, "templates"),
	}
	s.server.Handler = s.mux
	s.mux.HandleFunc("/", s.index)
	s.mux.PathPrefix("/static").Handler(
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
