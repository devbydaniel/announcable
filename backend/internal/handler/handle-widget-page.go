package handler

import (
	"errors"
	"html"
	"net/http"

	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type widgetPageData struct {
	BaseTemplateData
	Cfg                    *widgetconfigs.WidgetConfig
	SafeTitle              string
	SafeDescription        string
	SafeReleaseNoteCtaText string
}

var widgetTmpl = templates.Construct(
	"widget",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/widget.html",
)

func (h *Handler) HandleWidgetPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleWidgetPage")
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))
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

	// get widget config
	var cfg *widgetconfigs.WidgetConfig
	cfg, err := widgetService.Get(uuid.MustParse(orgId))
	if err != nil {
		if errors.Is(err, h.DB.ErrRecordNotFound) {
			// this should not happen, just in case...
			h.log.Warn().Msg("Widget config not found, creating...")
			cfg, err = widgetService.Init(uuid.MustParse(orgId))
			if err != nil {
				h.log.Error().Err(err).Msg("Error creating widget config")
				http.Error(w, "Error creating widget config", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		}
	}

	data := widgetPageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Widget Config",
			HasActiveSubscription: hasActiveSubscription,
		},
		Cfg:                    cfg,
		SafeTitle:              html.EscapeString(cfg.Title),
		SafeDescription:        html.EscapeString(cfg.Description),
		SafeReleaseNoteCtaText: html.EscapeString(cfg.ReleaseNoteCtaText),
	}
	if err := widgetTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
