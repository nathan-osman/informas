package server

import (
	"net/http"
)

// logout ends the user's current session.
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessions.Get(r, sessionName)
	delete(session.Values, sessionUserID)
	session.Save(r, w)
	s.addAlert(w, r, alertInfo, "you have been logged out")
	http.Redirect(w, r, "/login", http.StatusFound)
}
