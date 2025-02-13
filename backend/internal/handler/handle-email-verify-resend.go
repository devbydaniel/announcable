package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
)

type resendForm struct {
	Email string
}

func (h *Handler) HandleEmailVerifyResend(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleEmailVerifyResend")
	userService := user.NewService(*user.NewRepository(h.DB))
	sessionService := session.NewService(*session.NewRepository(h.DB))

	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := registerForm{
		Email: r.FormValue("email"),
	}
	h.log.Debug().Interface("req", req).Msg("Email resend request")

	user, err := userService.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, h.DB.ErrRecordNotFound) {
			h.log.Warn().Str("email", req.Email).Msg("User not found for resending email")
			w.Header().Set("HX-Trigger", "custom:submit-success")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.log.Error().Err(err).Msg("Error accessing user")
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	sessionService.InvalidateUserSessions(user.ID)

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
