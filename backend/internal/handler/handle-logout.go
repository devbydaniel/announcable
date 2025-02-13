package handler

import (
	"net/http"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/google/uuid"
)

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleLogout")
	ctx := r.Context()
	sessionId, ok := ctx.Value(mw.SessionIdKey).(string)
	if !ok {
		http.Error(w, "Error getting user id", http.StatusInternalServerError)
		return
	}

	sessionService := session.NewService(*session.NewRepository(h.DB))
	if err := sessionService.Delete(uuid.MustParse(sessionId)); err != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     session.AuthCookieName,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
