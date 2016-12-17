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

// view ensures that some basic sanity checks are run before processing the
// request. This includes permission checks and form pre-processing.
func (s *Server) view(a access, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *db.User
		session, _ := s.sessions.Get(r, sessionName)
		if v, ok := session.Values[sessionUserID]; ok {
			u, err := db.FindUser("ID", v.(int))
			if err == nil {
				user = u
			}
		}
		context.Set(r, contextUser, user)
		if a != accessPublic && user == nil ||
			a == accessRegistered && !user.IsAdmin {
			s.addAlert(w, r, alertDanger, "page requires authorization")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
		fn(w, r)
	}
}
