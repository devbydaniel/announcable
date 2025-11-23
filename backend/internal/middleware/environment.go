package mw

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
)

// CloudOnly middleware prevents access to routes in self-hosted mode
func (h *Handler) CloudOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.New()
		if cfg.IsSelfHosted() {
			http.Error(w, "This feature is not available in self-hosted mode", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
