package server

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/informas/db"
)

const (
	configInstalled = "installed"
)

// install is run for the first time when the "install" config option is unset.
// It is used to initialize the database and configure the first admin account.
func (s *Server) install(w http.ResponseWriter, r *http.Request) {
	var (
		adminUsername string = "admin"
		adminPassword string
		adminEmail    string
	)
	if r.Method == http.MethodPost {
		adminUsername = r.Form.Get("admin_username")
		adminPassword = r.Form.Get("admin_password")
		adminEmail = r.Form.Get("admin_email")
		u := &db.User{
			Username: adminUsername,
			Email:    adminEmail,
			IsAdmin:  true,
		}
		if err := u.SetPassword(adminPassword); err != nil {
			s.addAlert(w, r, alertDanger, err.Error())
		} else {
			if err := u.Save(); err != nil {
				s.addAlert(w, r, alertDanger, err.Error())
			} else {
				s.addAlert(w, r, alertInfo, "installation complete")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "install.html", pongo2.Context{
		"admin_username": adminUsername,
		"admin_password": adminPassword,
		"admin_email":    adminEmail,
	})
}
