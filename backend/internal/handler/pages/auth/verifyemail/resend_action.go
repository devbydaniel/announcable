package verifyemail

import (
	"errors"
	"net/http"
	"time"

	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
)

type resendForm struct {
	Email string
}

// HandleResend handles POST /verify-email/
func (h *Handlers) HandleResend(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleResend")
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))

	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := resendForm{
		Email: r.FormValue("email"),
	}
	h.deps.Log.Debug().Interface("req", req).Msg("Email resend request")

	user, err := userService.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
			h.deps.Log.Warn().Str("email", req.Email).Msg("User not found for resending email")
			w.Header().Set("HX-Trigger", "custom:submit-success")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.deps.Log.Error().Err(err).Msg("Error accessing user")
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	_ = sessionService.InvalidateUserSessions(user.ID)

	token := sessionService.CreateToken()
	if err := sessionService.CreateCustomDuration(token, user.ID, 1*time.Hour); err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	if err := userService.SendVerifcationEmail(user, token); err != nil {
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusCreated)
}
