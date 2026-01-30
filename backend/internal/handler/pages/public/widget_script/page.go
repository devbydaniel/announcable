package widget_script

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/static"
)

// Handlers provides widget script handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new widget script handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// ServeWidgetScript serves the Lit widget JavaScript file
func (h *Handlers) ServeWidgetScript(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeWidgetScript")
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(static.Widget)
}
