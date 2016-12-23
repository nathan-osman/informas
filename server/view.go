package server

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/nathan-osman/informas/db"
)

// access determines which users may view a page.
type access int

const (
	accessPublic access = iota
	accessRegistered
	accessAdmin
)

const contextUser = "user"

// view wraps each of the individual view functions. It takes care of such
// things as authentication, form processing, errors, etc.
func (s *Server) view(a access, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Always redirect to the installer if not installed
		if s.config.GetInt(configInstalled) == 0 && r.RequestURI != "/install" {
			http.Redirect(w, r, "/install", http.StatusFound)
			return
		}

		// Check for a user session and add it to the context if one exists
		var user *db.User
		session, _ := s.sessions.Get(r, sessionName)
		if v, ok := session.Values[sessionUserID]; ok {
			u, err := db.FindUser(&db.Token{}, "ID", v.(int))
			if err == nil {
				user = u
			}
		}
		context.Set(r, contextUser, user)

		// Confirm that the user has permission to access the view
		if a != accessPublic && user == nil || a == accessRegistered && !user.IsAdmin {
			s.addAlert(w, r, alertDanger, "page requires authorization")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// For POST requests, parse the form
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		// Execute the view
		f(w, r)
	}
}
