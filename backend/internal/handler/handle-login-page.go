package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/templates"
)

type loginPageData struct {
	BaseTemplateData
}

var loginTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/login.html",
)

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{
		BaseTemplateData: BaseTemplateData{
			Title: "Login",
		},
	}
	if err := loginTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
