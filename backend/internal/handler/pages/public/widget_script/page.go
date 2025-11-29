package widget_script

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
	"github.com/devbydaniel/release-notes-go/static"
)

// Handlers provides widget script handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new widget script handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// ServeWidgetScript serves the React widget JavaScript file
func (h *Handlers) ServeWidgetScript(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeWidgetScript")
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(static.Widget)
}

// ServeWidgetLitScript serves the Lit widget JavaScript file
func (h *Handlers) ServeWidgetLitScript(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeWidgetLitScript")
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(static.WidgetLit)
}
