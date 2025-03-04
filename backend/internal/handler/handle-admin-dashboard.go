package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/admin"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type adminDashboardData struct {
	BaseTemplateData
	Organisations []*OrganisationData
}

type OrganisationData struct {
	ID        string
	Name      string
	CreatedAt string
}

var adminDashboardTmpl = templates.Construct(
	"admin-dashboard",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/admin-dashboard.html",
)

func (h *Handler) HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleAdminDashboard")
	adminService := admin.NewService(*admin.NewRepository(h.DB))

	userId, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		h.log.Error().Msg("Error finding user")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if the user is an admin
	if !adminService.IsAdminUser(uuid.MustParse(userId)) {
		h.log.Warn().Str("userId", userId).Msg("Unauthorized access attempt to admin dashboard")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get all organisations
	orgs, err := adminService.GetAllOrganisations(uuid.MustParse(userId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting organisations")
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

	data := adminDashboardData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Admin Dashboard",
			HasActiveSubscription: true,
		},
		Organisations: orgData,
	}

	// Render the template
	if err := adminDashboardTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
