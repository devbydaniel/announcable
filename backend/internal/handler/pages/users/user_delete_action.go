package users

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandleUserDelete handles DELETE /users/{id}
func (h *Handlers) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleUserDelete")
	orgUserId := chi.URLParam(r, "id")
	h.deps.Log.Debug().Str("ouId", orgUserId).Msg("id URL param")
	if orgUserId == "" {
		h.deps.Log.Error().Msg("User ID not found in URL")
		http.Error(w, "Error deleting user", http.StatusBadRequest)
		return
	}

	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))

	ou, err := orgService.GetOrgUser(uuid.MustParse(orgUserId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting org user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := orgService.RemoveFromOrg(uuid.MustParse(orgUserId)); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error deleting user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := userService.Delete(ou.UserID); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error deleting user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if err := sessionService.InvalidateUserSessions(ou.UserID); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error deleting user session")
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
