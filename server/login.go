package server

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/informas/db"
)

// login presents the login form.
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		password string
	)
	if r.Method == http.MethodPost {
		username = r.Form.Get("username")
		password = r.Form.Get("password")
		u, err := db.FindUser("Username", username)
		if err != nil {
			s.addAlert(w, r, alertDanger, "invalid username")
		} else {
			if err := u.Authenticate(password); err != nil {
				s.addAlert(w, r, alertDanger, "invalid password")
			} else {
				session, _ := s.sessions.Get(r, sessionName)
				session.Values[sessionUserID] = u.ID
				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "login.html", pongo2.Context{
		"username": username,
		"password": password,
	})
}
