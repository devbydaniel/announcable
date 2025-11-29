package dashboard

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/admin"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/google/uuid"
)

// Handlers provides admin dashboard handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new admin dashboard handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// DashboardData represents the admin dashboard template data
type DashboardData struct {
	shared.BaseTemplateData
	Organisations []*OrganisationData
}

// OrganisationData represents organisation info for the dashboard
type OrganisationData struct {
	ID        string
	Name      string
	CreatedAt string
}

var dashboardTmpl = templates.Construct(
	"admin-dashboard",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/admin-dashboard.html",
)

// ServeDashboardPage renders the admin dashboard
func (h *Handlers) ServeDashboardPage(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeDashboardPage")
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

	// Get all organisations
	orgs, err := adminService.GetAllOrganisations(uuid.MustParse(userId))
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting organisations")
		http.Error(w, "Error getting organisations", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	orgData := make([]*OrganisationData, 0, len(orgs))
	for _, org := range orgs {
		orgData = append(orgData, &OrganisationData{
			ID:        org.ID.String(),
			Name:      org.Name,
			CreatedAt: org.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	data := DashboardData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Admin Dashboard",
		},
		Organisations: orgData,
	}

	// Render the template
	if err := dashboardTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.Log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
