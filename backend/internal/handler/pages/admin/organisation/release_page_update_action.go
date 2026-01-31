package organisation

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/admin"
	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type releasePageUpdateForm struct {
	Slug string `schema:"slug" validate:"required"`
}

// HandleReleasePageUpdate updates an organisation's release page slug
func (h *Handlers) HandleReleasePageUpdate(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleReleasePageUpdate")

	// Get the current user from the session
	adminService := admin.NewService(*admin.NewRepository(h.DB))
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))

	userID, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		h.Log.Error().Msg("Error finding user")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// Check if the user is an admin
	if !adminService.IsAdminUser(uuid.MustParse(userID)) {
		h.Log.Warn().Str("userID", userID).Msg("Unauthorized access attempt to admin functionality")
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get organisation ID from URL params
	orgIDStr := chi.URLParam(r, "orgID")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		h.Log.Error().Str("orgID", orgIDStr).Err(err).Msg("Error parsing organisation ID")
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
	var updateDTO releasePageUpdateForm
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

	// Update release page slug
	if err := releasePageConfigService.EditSlugAsAdmin(orgID, updateDTO.Slug); err != nil {
		h.Log.Error().Err(err).Msg("Error updating release page slug")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
