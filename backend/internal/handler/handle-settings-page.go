package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type settingsPageData struct {
	BaseTemplateData
	WidgetID            string
	ReleasePageUrl      string
	CustomUrl           *string
	HasPaidSubscription bool
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

	hasActiveSubscription, ok := ctx.Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))
	subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))

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
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release page URL")
	}

	isFree := false
	if hasActiveSubscription {
		var err error
		isFree, err = subscriptionService.IsFreeSubscription(uuid.MustParse(orgId))
		if err != nil {
			h.log.Error().Err(err).Msg("Error checking if subscription is free")
			http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
			return
		}
	}

	orgName := ctx.Value(mw.OrgNameKey).(string)

	pageData := settingsPageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Settings for " + orgName,
			HasActiveSubscription: hasActiveSubscription,
		},
		WidgetID:            externalId.String(),
		HasPaidSubscription: hasActiveSubscription && !isFree,
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
