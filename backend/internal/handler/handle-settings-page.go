package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type settingsPageData struct {
	Title          string
	WidgetID       string
	ReleasePageUrl string
	CustomUrl      *string
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
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))

	widgetConfig, err := widgetConfigService.Get(uuid.MustParse(orgId))
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

	var releasePageUrl string
	releasePageUrl, err = releasePageConfigService.GetUrl(uuid.MustParse(orgId))
	h.log.Debug().Str("releasePageUrl", releasePageUrl).Msg("Got release page URL")
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release page URL")
	} else if releasePageUrl == "" {
		orgName := ctx.Value(mw.OrgNameKey).(string)
		if err := releasePageConfigService.UpdateSlug(uuid.MustParse(orgId), orgName); err != nil {
			h.log.Error().Err(err).Msg("Error updating release page slug")
		} else {
			releasePageUrl, err = releasePageConfigService.GetUrl(uuid.MustParse(orgId))
			if err != nil {
				h.log.Error().Err(err).Msg("Error getting release page URL")
			}
		}
	}

	pageData := settingsPageData{
		Title:    "Settings",
		WidgetID: externalId.String(),
	}
	if releasePageUrl != "" {
		pageData.ReleasePageUrl = releasePageUrl
	}

	if widgetConfig.ReleasePageBaseUrl != nil {
		pageData.CustomUrl = widgetConfig.ReleasePageBaseUrl
	}

	if err := settingsTmpl.ExecuteTemplate(w, "root", pageData); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
