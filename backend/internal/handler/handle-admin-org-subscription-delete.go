package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/admin"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) HandleAdminOrgSubscriptionDelete(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleAdminOrgSubscriptionDelete")
	
	// Get the current user from the session
	adminService := admin.NewService(*admin.NewRepository(h.DB))

	userId, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		h.log.Error().Msg("Error finding user")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// Check if the user is an admin
	if !adminService.IsAdminUser(uuid.MustParse(userId)) {
		h.log.Warn().Str("userId", userId).Msg("Unauthorized access attempt to admin functionality")
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get subscription ID from URL
	subscriptionId := chi.URLParam(r, "id")
	if subscriptionId == "" {
		h.log.Error().Msg("Subscription ID not found in URL")
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	// Delete the subscription
	if err := adminService.DeleteSubscription(uuid.MustParse(userId), uuid.MustParse(subscriptionId)); err != nil {
		h.log.Error().Err(err).Msg("Error deleting subscription")
		http.Error(w, "Error deleting subscription", http.StatusInternalServerError)
		return
	}

	// Return empty response with OK status for HTMX to remove the row
	h.log.Info().Msg("Subscription deleted successfully")
	r.Header.Set("HX-Trigger", "custom:success")
	w.WriteHeader(http.StatusOK)
	return
}