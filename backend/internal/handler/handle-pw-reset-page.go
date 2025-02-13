package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
)

type PwResetPage struct {
	Title string
	Token string
}

var pwResetPageTmpl = templates.Construct(
	"invite-accept",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/reset-pw.html",
)

func (h *Handler) HandlePwResetPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandlePwResetPage")
	token := chi.URLParam(r, "token")

	data := InviteAcceptPage{
		Title: "Reset Password",
		Token: token,
	}
	if err := pwResetPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
