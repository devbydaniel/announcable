package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/static"
)

func (h *Handler) HandleWidgetjsServe(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleWidgetjsServe")
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(static.Widget)
}
