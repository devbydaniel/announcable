package handler

import (
	"encoding/json"
	"net/http"

	releasenotelikes "github.com/devbydaniel/release-notes-go/internal/domain/release-note-likes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type getLikeStateResponse struct {
	IsLiked bool `json:"is_liked"`
}

func (h *Handler) HandleGetReleaseNoteLikeState(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleGetReleaseNoteLikeState")

	// Get external org ID from URL params
	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.log.Error().Msg("Organisation ID not found in URL")
		http.Error(w, "Organisation ID required", http.StatusBadRequest)
		return
	}

	// Get release note ID from URL params
	releaseNoteId := chi.URLParam(r, "releaseNoteId")
	if releaseNoteId == "" {
		h.log.Error().Msg("Release note ID not found in URL")
		http.Error(w, "Release note ID required", http.StatusBadRequest)
		return
	}

	// Get client ID from URL params
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		// don't error, just return false
		fallback := getLikeStateResponse{
			IsLiked: false,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fallback); err != nil {
			h.log.Error().Err(err).Msg("Error encoding response")
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
		return
	}
	// Parse UUIDs
	releaseNoteUUID, err := uuid.Parse(releaseNoteId)
	if err != nil {
		h.log.Error().Err(err).Msg("Invalid release note ID")
		http.Error(w, "Invalid release note ID", http.StatusBadRequest)
		return
	}

	// Get like state
	likesService := releasenotelikes.NewService(releasenotelikes.NewRepository(h.DB))
	isLiked, err := likesService.HasUserLiked(releaseNoteUUID, clientId)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting like state")
		http.Error(w, "Error getting like state", http.StatusInternalServerError)
		return
	}

	// Return response
	response := getLikeStateResponse{
		IsLiked: isLiked,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Err(err).Msg("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
} 