package logout

import (
	"net/http"
	"time"

	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/google/uuid"
)

// Handlers holds the dependencies for logout handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// HandleLogout handles GET /logout/
func (h *Handlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleLogout")
	ctx := r.Context()
	sessionId, ok := ctx.Value(mw.SessionIdKey).(string)
	if !ok {
		http.Error(w, "Error getting user id", http.StatusInternalServerError)
		return
	}

	sessionService := session.NewService(*session.NewRepository(h.deps.DB))
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
