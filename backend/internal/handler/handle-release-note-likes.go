package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotelikes "github.com/devbydaniel/release-notes-go/internal/domain/release-note-likes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type toggleLikeRequest struct {
	ReleaseNoteID string `json:"release_note_id"`
	ClientID      string `json:"client_id"`
}

type toggleLikeResponse struct {
	IsLiked bool `json:"is_liked"`
	Count   int  `json:"count"`
}

func (h *Handler) HandleReleaseNoteToggleLike(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNoteToggleLike")
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))

	// Get external org ID from URL params
	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.log.Error().Msg("Organisation ID not found in URL")
		http.Error(w, "Organisation ID required", http.StatusBadRequest)
		return
	}

	// Get org by external ID
	org, err := orgService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting organisation")
		http.Error(w, "Error getting organisation", http.StatusInternalServerError)
		return
	}

	// Parse request body
	var req toggleLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Err(err).Msg("Error decoding request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ReleaseNoteID == "" || req.ClientID == "" {
		h.log.Error().Msg("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Parse UUIDs
	releaseNoteUUID, err := uuid.Parse(req.ReleaseNoteID)
	if err != nil {
		h.log.Error().Err(err).Msg("Invalid release note ID")
		http.Error(w, "Invalid release note ID", http.StatusBadRequest)
		return
	}

	orgUUID, err := uuid.Parse(org.ID.String())
	if err != nil {
		h.log.Error().Err(err).Msg("Invalid organisation ID")
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	// Toggle like
	likesService := releasenotelikes.NewService(releasenotelikes.NewRepository(h.DB))
	isLiked, err := likesService.ToggleLike(releaseNoteUUID, orgUUID, req.ClientID)
	if err != nil {
		h.log.Error().Err(err).Msg("Error toggling like")
		http.Error(w, "Error toggling like", http.StatusInternalServerError)
		return
	}

	// Return response with updated state
	response := toggleLikeResponse{
		IsLiked: isLiked,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Err(err).Msg("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
} 