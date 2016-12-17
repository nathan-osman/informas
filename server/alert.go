package server

import (
	"encoding/gob"
	"net/http"
)

const (
	alertInfo   = "info"
	alertDanger = "danger"
)

type alert struct {
	Type string
	Body string
}

func init() {
	gob.Register(&alert{})
}

// addAlert stores the provided alert in the current session for retrieval when
// the next page is rendered.
func (s *Server) addAlert(w http.ResponseWriter, r *http.Request, alertType, body string) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.AddFlash(&alert{
		Type: alertType,
		Body: body,
	})
}

// getAlerts retrieves all of the alerts for the current session.
func (s *Server) getAlerts(w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	return session.Flashes()
}
