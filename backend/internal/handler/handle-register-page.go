package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/templates"
)

type registerPageData struct {
	BaseTemplateData
}

var registerTmpl = templates.Construct(
	"register",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/register.html",
)

func (h *Handler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	data := registerPageData{
		BaseTemplateData: BaseTemplateData{
			Title: "Register",
		},
	}
	if err := registerTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
