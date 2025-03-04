package handler

import (
	"net/http"

	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
)

type PwResetPage struct {
	BaseTemplateData
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
	hasActiveSubscription, ok := r.Context().Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	data := PwResetPage{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Reset Password",
			HasActiveSubscription: hasActiveSubscription,
		},
		Token: token,
	}
	if err := pwResetPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
