package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/devbydaniel/announcable/templates"
)

var notFoundTmpl = templates.Construct(
	"not-found",
	"layouts/root.html",
	"layouts/fullscreenmessage.html",
	"pages/not-found.html",
)

type notFoundErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// HandleNotFound handles 404 errors
func (h *Handlers) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Str("path", r.URL.Path).Str("method", r.Method).Msg("HandleNotFound")
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		// JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		response := notFoundErrorResponse{
			Status:  404,
			Message: "Page not found",
			Path:    r.URL.Path,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		// HTML response
		w.WriteHeader(404)
		data := map[string]interface{}{
			"Title": "Page Not Found",
		}
		if err := notFoundTmpl.ExecuteTemplate(w, "root", data); err != nil {
			fmt.Fprintf(os.Stderr, "HandleNotFound: Error executing template: %v\n", err)
			h.Log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
