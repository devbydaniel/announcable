package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type settingsPageData struct {
	Title     string
	WidgetID  string
	CustomUrl *string
}

var settingsTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/settings-page.html",
)

func (h *Handler) HandleSettingsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))

	widgetConfig, err := widgetConfigService.Get(orgId)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	externalId, err := organisationService.GetExternalId(uuid.MustParse(orgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting external org ID")
		http.Error(w, "Error getting external org ID", http.StatusInternalServerError)
		return
	}

	pageData := settingsPageData{
		Title:    "Settings",
		WidgetID: externalId.String(),
	}

	if widgetConfig.ReleasePageBaseUrl != nil {
		pageData.CustomUrl = widgetConfig.ReleasePageBaseUrl
	}

	if err := settingsTmpl.ExecuteTemplate(w, "root", pageData); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
