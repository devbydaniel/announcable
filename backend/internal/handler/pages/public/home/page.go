package home

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
)

// Handlers provides home page handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new home page handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// ServeHomePage redirects to release notes page
func (h *Handlers) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/release-notes", http.StatusSeeOther)
}
