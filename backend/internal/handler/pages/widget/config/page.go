package config

import (
	"errors"
	"html"
	"net/http"

	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/google/uuid"
)

// Handlers holds dependencies for widget config handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData represents the template data for widget config
type pageData struct {
	shared.BaseTemplateData
	Cfg                    *widgetconfigs.WidgetConfig
	SafeTitle              string
	SafeDescription        string
	SafeReleaseNoteCtaText string
}

var pageTmpl = templates.Construct(
	"widget",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/widget.html",
)

// ServeWidgetConfigPage handles GET /widget-config/
func (h *Handlers) ServeWidgetConfigPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeWidgetConfigPage")
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))

	orgID, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// get widget config
	var cfg *widgetconfigs.WidgetConfig
	cfg, err := widgetService.Get(uuid.MustParse(orgID))
	if err != nil {
		if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
			// this should not happen, just in case...
			h.deps.Log.Warn().Msg("Widget config not found, creating...")
			cfg, err = widgetService.Init(uuid.MustParse(orgID))
			if err != nil {
				h.deps.Log.Error().Err(err).Msg("Error creating widget config")
				http.Error(w, "Error creating widget config", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Error getting widget config", http.StatusInternalServerError)
			return
		}
	}

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Widget Config",
		},
		Cfg:                    cfg,
		SafeTitle:              html.EscapeString(cfg.Title),
		SafeDescription:        html.EscapeString(cfg.Description),
		SafeReleaseNoteCtaText: html.EscapeString(cfg.ReleaseNoteCtaText),
	}

	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
