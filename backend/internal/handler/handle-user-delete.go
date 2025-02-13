package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleUserDelete")
	orgUserId := chi.URLParam(r, "id")
	h.log.Debug().Str("ouId", orgUserId).Msg("id URL param")
	if orgUserId == "" {
		h.log.Error().Msg("User ID not found in URL")
		http.Error(w, "Error deleting user", http.StatusBadRequest)
		return
	}

	orgService := organisation.NewService(*organisation.NewRepository(h.DB))
	userService := user.NewService(*user.NewRepository(h.DB))
	sessionService := session.NewService(*session.NewRepository(h.DB))

	ou, err := orgService.GetOrgUser(uuid.MustParse(orgUserId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting org user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := orgService.RemoveFromOrg(uuid.MustParse(orgUserId)); err != nil {
		h.log.Error().Err(err).Msg("Error deleting user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := userService.Delete(ou.UserID); err != nil {
		h.log.Error().Err(err).Msg("Error deleting user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := sessionService.InvalidateUserSessions(ou.UserID); err != nil {
		h.log.Error().Err(err).Msg("Error deleting user session")
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
	return
}
