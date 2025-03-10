package handler

import (
	"net/http"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/internal/ratelimit"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-playground/validator"
)

type pwForgotForm struct {
	Email string `json:"email" validate:"required,email"`
}

var pwForgotTmpl = templates.Construct("pw-reset-confirm", "partials/hx-pw-forgot-confirm.html")
var pwForgotRateLimiter = ratelimit.New(60, 2)

func (h *Handler) HandlePwForgot(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandlePwForgot")
	userService := user.NewService(*user.NewRepository(h.DB))
	sessionService := session.NewService(*session.NewRepository(h.DB))

	// parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// decode form
	var forgotForm pwForgotForm
	if err := h.decoder.Decode(&forgotForm, r.PostForm); err != nil {
		h.log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(forgotForm); err != nil {
		h.log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// check rate limit
	if err := pwForgotRateLimiter.Deduct(forgotForm.Email, 1); err != nil {
		h.log.Warn().Str("email", forgotForm.Email).Msg("Rate limit exceeded for password reset requests")
		http.Error(w, "Too many password reset requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// check if user exists
	usr, err := userService.GetByEmail(forgotForm.Email)
	if err != nil {
		// don't reveal if user exists
		if err := pwForgotTmpl.ExecuteTemplate(w, "hx-pw-reset-confirm", nil); err != nil {
			h.log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}

	// invalidate session, create new session, send email
	sessionService.InvalidateUserSessions(usr.ID)
	token := sessionService.CreateToken()
	if err := sessionService.CreateCustomDuration(token, usr.ID, 1*time.Hour); err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	userService.SendPwResetEmail(usr, token)

	if err := pwForgotTmpl.ExecuteTemplate(w, "hx-pw-reset-confirm", nil); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
