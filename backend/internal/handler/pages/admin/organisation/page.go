package organisation

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/admin"
	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handlers provides admin organisation handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new admin organisation handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// OrganisationDetailsData represents the organisation details template data
type OrganisationDetailsData struct {
	shared.BaseTemplateData
	Organisation *OrganisationDetailData
	ReleasePage  *ReleasePageData
	Users        []*OrganisationUserData
}

// OrganisationDetailData represents organisation detail info
type OrganisationDetailData struct {
	ID        string
	Name      string
	CreatedAt string
}

// ReleasePageData represents release page info
type ReleasePageData struct {
	Slug string
}

// OrganisationUserData represents user info in an organisation
type OrganisationUserData struct {
	ID        string
	Email     string
	Role      string
	CreatedAt string
}

var orgDetailsTmpl = templates.Construct(
	"admin-org-details",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/admin-org-details.html",
)

// ServeOrganisationDetailsPage renders the organisation details page
func (h *Handlers) ServeOrganisationDetailsPage(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeOrganisationDetailsPage")
	orgId := chi.URLParam(r, "orgId")

	// Get the current user from the session
	adminService := admin.NewService(*admin.NewRepository(h.DB))

	userId, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		h.Log.Error().Msg("Error finding user")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if the user is an admin
	if !adminService.IsAdminUser(uuid.MustParse(userId)) {
		h.Log.Warn().Str("userId", userId).Msg("Unauthorized access attempt to admin dashboard")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get organisation details with users
	org, orgUsers, err := adminService.GetOrganisationWithUsers(uuid.MustParse(userId), uuid.MustParse(orgId))
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting organisation details")
		http.Error(w, "Error getting organisation details", http.StatusInternalServerError)
		return
	}

	// Get release page details
	releasePageService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))
	releasePageConfig, err := releasePageService.Get(uuid.MustParse(orgId))
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting release page details")
		http.Error(w, "Error getting release page details", http.StatusInternalServerError)
		return
	}
	h.Log.Debug().Interface("releasePageConfig", releasePageConfig).Msg("releasePageConfig")

	// Prepare data for the template
	orgData := &OrganisationDetailData{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	releasePageData := &ReleasePageData{
		Slug: releasePageConfig.Slug,
	}

	userData := make([]*OrganisationUserData, 0, len(orgUsers))
	for _, ou := range orgUsers {
		userData = append(userData, &OrganisationUserData{
			ID:        ou.User.ID.String(),
			Email:     ou.User.Email,
			Role:      ou.Role.String(),
			CreatedAt: ou.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	data := OrganisationDetailsData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: org.Name,
		},
		Organisation: orgData,
		ReleasePage:  releasePageData,
		Users:        userData,
	}

	// Render the template
	if err := orgDetailsTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.Log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
