package server

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/informas/db"
)

// settings allow site-wide configuration to be edited.
func (s *Server) settings(w http.ResponseWriter, r *http.Request) {
	var (
		siteTitle = s.config.GetString(configSiteTitle)
	)
	if r.Method == http.MethodPost {
		siteTitle = r.Form.Get("site_title")
		err := db.Transaction(func(t *db.Token) error {
			if err := s.config.SetString(t, configSiteTitle, siteTitle); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			s.addAlert(w, r, alertDanger, err.Error())
		} else {
			s.addAlert(w, r, alertInfo, "settings saved")
			http.Redirect(w, r, "/settings", http.StatusFound)
			return
		}
	}
	s.render(w, r, "settings.html", pongo2.Context{
		"title":       "Settings",
		"site_title_": siteTitle,
	})
}
