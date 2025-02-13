package handler

import (
	"net/http"

	lpconfigs "github.com/devbydaniel/release-notes-go/internal/domain/landing-page-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

type landingPageData struct {
	Title string
	Cfg   *lpconfigs.LpConfig
}

var landingPageTmpl = templates.Construct(
	"website",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/landing-page-config.html",
)

func (h *Handler) HandleLpConfigPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleLandingPage")
	lpService := lpconfigs.NewService(*lpconfigs.NewRepository(h.DB, h.ObjStore))
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// get landing page config
	cfg, err := lpService.Get(orgId)
	if err != nil {
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
	}

	data := landingPageData{
		Title: "Landing Page Config",
		Cfg:   cfg,
	}
	if err := landingPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
