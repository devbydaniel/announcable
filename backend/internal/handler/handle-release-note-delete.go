package handler

import (
	"fmt"
	"net/http"
	"net/url"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) HandleReleaseNoteDelete(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNoteDelete")
	rnId := chi.URLParam(r, "id")
	h.log.Debug().Str("rnId", rnId).Msg("id URL param")
	if rnId == "" {
		h.log.Error().Msg("Release note ID not found in URL")
		http.Error(w, "Error deleting release note", http.StatusBadRequest)
		return
	}

	releaseNoteService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
	if err := releaseNoteService.Delete(uuid.MustParse(rnId)); err != nil {
		h.log.Error().Err(err).Msg("Error deleting release note")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	successMsg := "release note deleted"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := fmt.Sprintf("/release-notes?success=%s", escapedMsg)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
	return
}
