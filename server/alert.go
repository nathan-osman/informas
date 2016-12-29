package server

import (
	"encoding/gob"
	"net/http"
)

// alertType determines what type of alert is shown.
type alertType string

const (
	alertInfo   alertType = "info"
	alertDanger alertType = "danger"
)

// alert is a message displayed to a user on the next page render.
type alert struct {
	Type alertType
	Body string
}

func init() {
	gob.Register(&alert{})
}

// addAlert stores the provided alert in the current session for retrieval when
// the next page is rendered.
func (s *Server) addAlert(w http.ResponseWriter, r *http.Request, type_ alertType, body string) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.AddFlash(&alert{
		Type: type_,
		Body: body,
	})
}

// getAlerts retrieves all of the alerts for the current session.
func (s *Server) getAlerts(w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	return session.Flashes()
}
