package account

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

// Handlers holds dependencies for settings/account handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData represents the template data for settings page
type pageData struct {
	shared.BaseTemplateData
	WidgetID            string
	ReleasePageUrl      string
	CustomUrl           *string
	HasPaidSubscription bool
}

var pageTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/settings-page.html",
)

// ServeSettingsPage handles GET /settings/
func (h *Handlers) ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	hasActiveSubscription, ok := ctx.Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.deps.Log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	organisationService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))
	subscriptionService := subscription.NewService(*subscription.NewRepository(h.deps.DB))

	widgetConfig, err := widgetConfigService.Get(uuid.MustParse(orgId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	externalId, err := organisationService.GetExternalId(uuid.MustParse(orgId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting external org ID")
		http.Error(w, "Error getting external org ID", http.StatusInternalServerError)
		return
	}

	var releasePageUrl string
	releasePageUrl, err = releasePageConfigService.GetUrl(uuid.MustParse(orgId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting release page URL")
	}

	orgName := ctx.Value(mw.OrgNameKey).(string)
	cfg := config.New()

	isFree := false
	// Only check subscription tier in cloud mode
	if cfg.IsCloud() && hasActiveSubscription {
		var err error
		isFree, err = subscriptionService.IsFreeSubscription(uuid.MustParse(orgId))
		if err != nil {
			h.deps.Log.Error().Err(err).Msg("Error checking if subscription is free")
			http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
			return
		}
	}

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title:                 "Settings for " + orgName,
			HasActiveSubscription: hasActiveSubscription,
			ShowSubscriptionUI:    cfg.IsCloud(),
		},
		WidgetID:            externalId.String(),
		HasPaidSubscription: hasActiveSubscription && !isFree,
	}
	if releasePageUrl != "" {
		data.ReleasePageUrl = releasePageUrl
	}

	if widgetConfig.ReleasePageBaseUrl != nil {
		data.CustomUrl = widgetConfig.ReleasePageBaseUrl
	}

	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
