package organisation

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/admin"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type orgUpdateForm struct {
	Name string `schema:"name" validate:"required"`
}

// HandleOrgUpdate updates an organisation's details
func (h *Handlers) HandleOrgUpdate(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleOrgUpdate")

	// Get the current user from the session
	adminService := admin.NewService(*admin.NewRepository(h.DB))
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))

	userId, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		h.Log.Error().Msg("Error finding user")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// Check if the user is an admin
	if !adminService.IsAdminUser(uuid.MustParse(userId)) {
		h.Log.Warn().Str("userId", userId).Msg("Unauthorized access attempt to admin functionality")
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get organisation ID from URL params
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		h.Log.Error().Str("orgId", orgIDStr).Err(err).Msg("Error parsing organisation ID")
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Decode form
	var updateDTO orgUpdateForm
	if err := h.Decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	// Validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	// Validate org name
	if err := orgService.IsValidOrgName(updateDTO.Name); err != nil {
		h.Log.Error().Err(err).Msg("Invalid organisation name")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update organisation
	if err := orgService.UpdateOrg(orgID, &organisation.Organisation{Name: updateDTO.Name}); err != nil {
		h.Log.Error().Err(err).Msg("Error updating organisation")
		http.Error(w, "Error updating organisation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
