package handler

import (
	"html"
	"net/http"

	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type landingPageData struct {
	Title             string
	Cfg               *releasepageconfig.ReleasePageConfig
	SafeTitle         string
	SafeDescription   string
	SafeBackLinkLabel string
}

var landingPageTmpl = templates.Construct(
	"website",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/landing-page-config.html",
)

func (h *Handler) HandleReleasePageConfigPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleasePageConfigPage")
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// get release page config
	cfg, err := releasePageConfigService.Get(uuid.MustParse(orgId))
	if err != nil {
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
	}

	data := landingPageData{
		Title:             "Release Page Config",
		Cfg:               cfg,
		SafeTitle:         html.EscapeString(cfg.Title),
		SafeDescription:   html.EscapeString(cfg.Description),
		SafeBackLinkLabel: html.EscapeString(cfg.BackLinkLabel),
	}
	if err := landingPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
