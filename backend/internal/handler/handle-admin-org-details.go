package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/admin"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type adminOrgDetailsData struct {
	BaseTemplateData
	Organisation  *OrganisationDetailData
	Users         []*OrganisationUserData
	Subscriptions []*SubscriptionData
}

type OrganisationDetailData struct {
	ID        string
	Name      string
	CreatedAt string
}

type OrganisationUserData struct {
	ID        string
	Email     string
	Role      string
	CreatedAt string
}

type SubscriptionData struct {
	ID        string
	CreatedAt          string
	IsActive           bool
	IsFree             bool
	StripeSubscriptionID string
}

var adminOrgDetailsTmpl = templates.Construct(
	"admin-org-details",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/admin-org-details.html",
)

func (h *Handler) HandleAdminOrgDetails(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleAdminOrgDetails")
	orgId := chi.URLParam(r, "orgId")

	// Get the current user from the session
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

	// Get organisation details with users
	org, orgUsers, err := adminService.GetOrganisationWithUsers(uuid.MustParse(userId), uuid.MustParse(orgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting organisation details")
		http.Error(w, "Error getting organisation details", http.StatusInternalServerError)
		return
	}

	// Get subscriptions
	subscriptions, err := adminService.GetSubscriptions(uuid.MustParse(userId), uuid.MustParse(orgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting subscriptions")
		http.Error(w, "Error getting subscriptions", http.StatusInternalServerError)
		return
	}
	h.log.Debug().Interface("subscriptions", subscriptions).Msg("subscriptions")

	// Prepare data for the template
	orgData := &OrganisationDetailData{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt.Format("2006-01-02 15:04:05"),
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

	subscriptionData := make([]*SubscriptionData, 0, len(subscriptions))
	for _, s := range subscriptions {
		subscriptionData = append(subscriptionData, &SubscriptionData{
			ID:                 s.ID.String(),
			CreatedAt:          s.CreatedAt.Format("2006-01-02 15:04:05"),
			IsActive:           s.IsActive,
			IsFree:             s.IsFree,
			StripeSubscriptionID: s.StripeSubscriptionID,
		})
	}

	data := adminOrgDetailsData{
		BaseTemplateData: BaseTemplateData{
			Title:                 org.Name,
			HasActiveSubscription: true,
		},
		Organisation:  orgData,
		Users:         userData,
		Subscriptions: subscriptionData,
	}

	// Render the template
	if err := adminOrgDetailsTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
