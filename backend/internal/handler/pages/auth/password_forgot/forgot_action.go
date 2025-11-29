package password_forgot

import (
	"net/http"
	"time"

	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
	"github.com/devbydaniel/announcable/internal/ratelimit"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-playground/validator"
)

type forgotForm struct {
	Email string `json:"email" validate:"required,email"`
}

var forgotTmpl = templates.Construct("pw-reset-confirm", "partials/hx-pw-forgot-confirm.html")
var forgotRateLimiter = ratelimit.New(60, 2)

// HandleForgotPassword handles POST /forgot-pw/
func (h *Handlers) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleForgotPassword")
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// decode form
	var forgotForm forgotForm
	if err := h.deps.Decoder.Decode(&forgotForm, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(forgotForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// check rate limit
	if err := forgotRateLimiter.Deduct(forgotForm.Email, 1); err != nil {
		h.deps.Log.Warn().Str("email", forgotForm.Email).Msg("Rate limit exceeded for password reset requests")
		http.Error(w, "Too many password reset requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// check if user exists
	usr, err := userService.GetByEmail(forgotForm.Email)
	if err != nil {
		// don't reveal if user exists
		if err := forgotTmpl.ExecuteTemplate(w, "hx-pw-reset-confirm", nil); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
		return
	}

	// invalidate session, create new session, send email
	sessionService.InvalidateUserSessions(usr.ID)
	token := sessionService.CreateToken()
	if err := sessionService.CreateCustomDuration(token, usr.ID, 1*time.Hour); err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	userService.SendPwResetEmail(usr, token)

	if err := forgotTmpl.ExecuteTemplate(w, "hx-pw-reset-confirm", nil); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
