package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) HandleInviteDelete(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleInviteDelete")
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))
	id := chi.URLParam(r, "id")
	if err := orgService.DeleteInvite(uuid.MustParse(id)); err != nil {
		h.log.Error().Err(err).Msg("Error deleting invite")
		http.Error(w, "Error deleting invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
	return
}
