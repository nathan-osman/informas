package server

import (
	"errors"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/informas/db"
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

// usersCreateOrEdit enables both new user accounts to be created and existing
// user accounts to be modified. Certain fields are only editable by
// administrators.
func (s *Server) usersCreateOrEdit(w http.ResponseWriter, r *http.Request, title, action string) {
	var (
		currentUser = context.Get(r, contextCurrentUser).(*db.User)
		user        = &db.User{}
		userID      = atoi(mux.Vars(r)["id"])
		password    = r.Form.Get("password")
		password2   = r.Form.Get("password2")
	)
	if !currentUser.IsAdmin && currentUser.ID != userID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err := db.Transaction(func(t *db.Token) error {
		if action == "edit" {
			u, err := db.FindUser(t, "ID", userID)
			if err != nil {
				return errors.New("invalid user")
			}
			user = u
		}
		if r.Method == http.MethodPost {
			if action == "edit" {
				if len(password) != 0 {
					if password != password2 {
						return errors.New("passwords do not match")
					}
					if err := user.SetPassword(password); err != nil {
						return errors.New("unable to set password")
					}
				}
			}
			user.Username = r.Form.Get("username")
			user.Email = r.Form.Get("email")
			if currentUser.IsAdmin {
				user.IsAdmin = len(r.Form.Get("is_admin")) != 0
				user.IsDisabled = len(r.Form.Get("is_disabled")) != 0
			}
			if err := user.Save(t); err != nil {
				return errors.New("unable to save user")
			}
		}
		return nil
	})
	if err != nil {
		s.addAlert(w, r, alertDanger, err.Error())
	} else if r.Method == http.MethodPost {
		s.addAlert(w, r, alertInfo, "user account saved")
		if currentUser.IsAdmin {
			http.Redirect(w, r, "/users", http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		return
	}
	s.render(w, r, "usersCreateOrEdit.html", pongo2.Context{
		"title":     title,
		"action":    action,
		"user":      user,
		"password":  password,
		"password2": password2,
	})
}

func (s *Server) usersCreate(w http.ResponseWriter, r *http.Request) {
	s.usersCreateOrEdit(w, r, "Create User", "create")
}

// usersIdEdit allows existing users to be modified.
func (s *Server) usersIdEdit(w http.ResponseWriter, r *http.Request) {
	s.usersCreateOrEdit(w, r, "Edit User", "edit")
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
