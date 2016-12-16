package server

import (
	"github.com/gorilla/mux"
	"github.com/hectane/go-asyncserver"
)

// Server provides the web interface for the application.
type Server struct {
	server *server.AsyncServer
	mux    *mux.Router
}

// New creates a new server instance.
func New(addr string) (*Server, error) {
	s := &Server{
		server: server.New(addr),
		mux:    mux.NewRouter(),
	}
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
