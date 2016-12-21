package server

import (
	"errors"
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
		err := db.Transaction(func(t *db.Token) error {
			username = r.Form.Get("username")
			password = r.Form.Get("password")
			u, err := db.FindUser(t, "Username", username)
			if err != nil {
				return errors.New("invalid username")
			}
			if err := u.Authenticate(password); err != nil {
				return errors.New("invalid password")
			}
			session, _ := s.sessions.Get(r, sessionName)
			session.Values[sessionUserID] = u.ID
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		})
		if err != nil {
			s.addAlert(w, r, alertDanger, err.Error())
		} else {
			return
		}
	}
	s.render(w, r, "login.html", pongo2.Context{
		"username": username,
		"password": password,
	})
}
