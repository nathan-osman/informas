package server

import (
	"errors"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/informas/db"
)

const (
	configInstalled = "installed"
	configSiteTitle = "site_title"
)

// install is used to initialize the application.
func (s *Server) install(w http.ResponseWriter, r *http.Request) {
	var (
		adminUsername string = "admin"
		adminPassword string
		adminEmail    string
	)
	if r.Method == http.MethodPost {
		err := db.Transaction(func(t *db.Token) error {

			// Do a basic sanity check on the form input
			adminUsername = r.Form.Get("admin_username")
			adminPassword = r.Form.Get("admin_password")
			adminEmail = r.Form.Get("admin_email")
			if adminUsername == "" || adminPassword == "" || adminEmail == "" {
				return errors.New("all fields are required")
			}

			// Create the initial admin
			u := &db.User{
				Username: adminUsername,
				Email:    adminEmail,
				IsAdmin:  true,
			}
			if err := u.SetPassword(adminPassword); err != nil {
				return err
			}
			if err := u.Save(t); err != nil {
				return err
			}

			// Create the initial configuration
			initialConfig := map[string]string{
				configInstalled: "1",
				configSiteTitle: "Informas",
			}
			for k, v := range initialConfig {
				if err := s.config.SetString(t, k, v); err != nil {
					return err
				}
			}

			// Indicate success
			s.addAlert(w, r, alertInfo, "installation complete")
			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		})
		if err != nil {
			s.addAlert(w, r, alertDanger, err.Error())
		} else {
			return
		}
	}
	s.render(w, r, "install.html", pongo2.Context{
		"title":          "Install",
		"admin_username": adminUsername,
		"admin_password": adminPassword,
		"admin_email":    adminEmail,
	})
}
