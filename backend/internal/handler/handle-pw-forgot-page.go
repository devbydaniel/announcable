package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/templates"
)

type forgotPwPageData struct {
	Title string
}

var forgotPwPageTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/forgot-pw.html",
)

func (h *Handler) HandlePwForgotPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandlePwForgotPage")
	data := loginPageData{
		Title: "Password reset",
	}
	if err := forgotPwPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
