package password_reset

import (
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/internal/password"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type resetForm struct {
	Password string `json:"password" validate:"required"`
	Confirm  string `json:"confirm" validate:"required"`
}

// HandleResetPassword handles POST /reset-pw/{token}/
func (h *Handlers) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleResetPassword")
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))

	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	h.deps.Log.Debug().Str("token", token).Msg("Token")

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// decode form
	var resetForm resetForm
	if err := h.deps.Decoder.Decode(&resetForm, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(resetForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	if resetForm.Password != resetForm.Confirm {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	if err := password.IsValidPassword(resetForm.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get user from token
	session, err := sessionService.ValidateSession(token)
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error validating session")
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}
	userId := session.UserID

	// invalidate session, create new session, send email
	if err := sessionService.InvalidateUserSessions(userId); err != nil {
		h.deps.Log.Warn().Err(err).Msg("Error invalidating session")
	}

	// update password
	if err := userService.UpdatePassword(userId, resetForm.Password); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating password")
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	successMsg := url.QueryEscape("Password updated")
	w.Header().Set("HX-Redirect", "/login?success="+successMsg)
	w.WriteHeader(http.StatusAccepted)
}
