package account

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
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
	WidgetID           string
	ReleasePageURL     string
	CustomURL          *string
	DisableReleasePage bool
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
	orgID, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	organisationService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))

	widgetConfig, err := widgetConfigService.Get(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	externalID, err := organisationService.GetExternalID(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting external org ID")
		http.Error(w, "Error getting external org ID", http.StatusInternalServerError)
		return
	}

	var releasePageURL string
	releasePageURL, err = releasePageConfigService.GetURL(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting release page URL")
	}

	releasePageConfig, err := releasePageConfigService.Get(uuid.MustParse(orgID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting release page config")
		http.Error(w, "Error getting release page config", http.StatusInternalServerError)
		return
	}

	orgName := ctx.Value(mw.OrgNameKey).(string)

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Settings for " + orgName,
		},
		WidgetID:           externalID.String(),
		DisableReleasePage: releasePageConfig.DisableReleasePage,
	}
	if releasePageURL != "" {
		data.ReleasePageURL = releasePageURL
	}

	if widgetConfig.ReleasePageBaseURL != nil {
		data.CustomURL = widgetConfig.ReleasePageBaseURL
	}

	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
