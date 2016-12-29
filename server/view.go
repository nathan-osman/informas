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

// view wraps each of the individual view functions. It ensures that
// installation has been completed, that the current user may access the page,
// and parses forms for POST requests.
func (s *Server) view(a access, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Always redirect to the installer if not installed
		if s.config.GetInt(configInstalled) == 0 && r.RequestURI != "/install" {
			http.Redirect(w, r, "/install", http.StatusFound)
			return
		}

		// Check for a user session
		var currentUser *db.User
		session, _ := s.sessions.Get(r, sessionName)
		if v, ok := session.Values[sessionUserID]; ok {
			u, err := db.FindUser(&db.Token{}, "ID", v.(int))
			if err == nil {
				currentUser = u
			}
		}
		context.Set(r, contextCurrentUser, currentUser)

		// Confirm that the user has permission to access the view
		if a != accessPublic && currentUser == nil || a == accessAdmin && !currentUser.IsAdmin {
			s.addAlert(w, r, alertDanger, "page requires authorization")
			http.Redirect(w, r, "/users/login", http.StatusFound)
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
