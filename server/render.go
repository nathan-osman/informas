package server

import (
	"net/http"
	"path"

	"github.com/flosch/pongo2"
)

const contextAlerts = "alerts"

// render loads the specified template, injects the provided context, and
// renders it directly to the response.
func (s *Server) render(w http.ResponseWriter, r *http.Request, templateName string, ctx pongo2.Context) {
	t, err := pongo2.FromFile(path.Join(s.templateDir, templateName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx[contextAlerts] = s.getAlerts(w, r)
	b, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
