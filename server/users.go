package server

import (
	"errors"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/informas/db"
)

const (
	sessionUserID = "user_id"
)

func (s *Server) usersIndex(w http.ResponseWriter, r *http.Request) {
	u, err := db.AllUsers(&db.Token{}, "Username")
	if err != nil {
		s.addAlert(w, r, alertDanger, err.Error())
	}
	s.render(w, r, "usersIndex.html", pongo2.Context{
		"title": "Users",
		"users": u,
	})
}

func (s *Server) usersId(w http.ResponseWriter, r *http.Request) {
	//...
}

func (s *Server) usersIdDelete(w http.ResponseWriter, r *http.Request) {
	//...
}

func (s *Server) usersLogin(w http.ResponseWriter, r *http.Request) {
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
			if u.IsDisabled {
				return errors.New("disabled account")
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
	s.render(w, r, "usersLogin.html", pongo2.Context{
		"title":    "Login",
		"username": username,
		"password": password,
	})
}

func (s *Server) usersLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessions.Get(r, sessionName)
	delete(session.Values, sessionUserID)
	session.Save(r, w)
	s.addAlert(w, r, alertInfo, "you have been logged out")
	http.Redirect(w, r, "/users/login", http.StatusFound)
}
