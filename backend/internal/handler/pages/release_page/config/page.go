package config

import (
	"errors"
	"html"
	"net/http"

	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/google/uuid"
)

// Handlers holds dependencies for release page config handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData represents the template data for release page config
type pageData struct {
	shared.BaseTemplateData
	Cfg               *releasepageconfig.ReleasePageConfig
	SafeTitle         string
	SafeDescription   string
	SafeBackLinkLabel string
	ReleasePageURL    string
}

var pageTmpl = templates.Construct(
	"website",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/landing-page-config.html",
)

// ServeReleasePageConfigPage handles GET /release-page-config/
func (h *Handlers) ServeReleasePageConfigPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeReleasePageConfigPage")
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))

	orgID, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// get release page config
	var cfg *releasepageconfig.ReleasePageConfig
	cfg, err := releasePageConfigService.Get(uuid.MustParse(orgID))
	if err != nil {
		if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
			// this should not happen, just in case...
			h.deps.Log.Warn().Msg("Release page config not found, creating...")
			orgName := r.Context().Value(mw.OrgNameKey).(string)
			cfg, err = releasePageConfigService.Init(uuid.MustParse(orgID), orgName)
			if err != nil {
				h.deps.Log.Error().Err(err).Msg("Error creating release page config")
				http.Error(w, "Error creating release page config", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Error getting widget config", http.StatusInternalServerError)
			return
		}
	}

	// get release page URL, either from slug or custom URL
	releasePageURL, err := releasePageConfigService.GetURL(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting release page URL")
	}
	widgetCfg, err := widgetConfigService.Get(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting widget config")
	}
	if widgetCfg.ReleasePageBaseURL != nil {
		releasePageURL = *widgetCfg.ReleasePageBaseURL
	}
	h.deps.Log.Debug().Str("releasePageURL", releasePageURL).Msg("Release page URL")

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Release Page Config",
		},
		Cfg:               cfg,
		SafeTitle:         html.EscapeString(cfg.Title),
		SafeDescription:   html.EscapeString(cfg.Description),
		SafeBackLinkLabel: html.EscapeString(cfg.BackLinkLabel),
		ReleasePageURL:    releasePageURL,
	}

	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
