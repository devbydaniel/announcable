package handler

import (
	"net/http"
)

func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/release-notes", http.StatusSeeOther)
}
