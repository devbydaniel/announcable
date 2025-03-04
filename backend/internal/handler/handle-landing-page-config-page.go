package handler

import (
	"errors"
	"html"
	"net/http"

	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type releasePageConfigPageData struct {
	BaseTemplateData
	Cfg               *releasepageconfig.ReleasePageConfig
	SafeTitle         string
	SafeDescription   string
	SafeBackLinkLabel string
}

var releasePageConfigPageTmpl = templates.Construct(
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
	hasActiveSubscription, ok := r.Context().Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}

	// get release page config
	var cfg *releasepageconfig.ReleasePageConfig
	cfg, err := releasePageConfigService.Get(uuid.MustParse(orgId))
	if err != nil {
		if errors.Is(err, h.DB.ErrRecordNotFound) {
			// this should not happen, just in case...
			h.log.Warn().Msg("Release page config not found, creating...")
			orgName := r.Context().Value(mw.OrgNameKey).(string)
			cfg, err = releasePageConfigService.Init(uuid.MustParse(orgId), orgName)
			if err != nil {
				h.log.Error().Err(err).Msg("Error creating release page config")
				http.Error(w, "Error creating release page config", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		}
	}

	data := releasePageConfigPageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Release Page Config",
			HasActiveSubscription: hasActiveSubscription,
		},
		Cfg:               cfg,
		SafeTitle:         html.EscapeString(cfg.Title),
		SafeDescription:   html.EscapeString(cfg.Description),
		SafeBackLinkLabel: html.EscapeString(cfg.BackLinkLabel),
	}
	if err := releasePageConfigPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
