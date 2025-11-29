package users

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandleInviteDelete handles DELETE /invites/{id}
func (h *Handlers) HandleInviteDelete(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleInviteDelete")
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	id := chi.URLParam(r, "id")
	if err := orgService.DeleteInvite(uuid.MustParse(id)); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error deleting invite")
		http.Error(w, "Error deleting invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
