package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/devbydaniel/release-notes-go/templates"
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

func (h *Handler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Str("path", r.URL.Path).Str("method", r.Method).Msg("HandleNotFound")
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
			h.log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
