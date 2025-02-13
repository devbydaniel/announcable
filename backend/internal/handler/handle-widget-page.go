package handler

import (
	"net/http"

	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

type widgetPageData struct {
	Title string
	Cfg   *widgetconfigs.WidgetConfig
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

	// get widget config
	cfg, err := widgetService.Get(orgId)
	if err != nil {
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
	}

	data := widgetPageData{
		Title: "Widget",
		Cfg:   cfg,
	}
	if err := widgetTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
