package handler

import (
	"errors"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/internal/password"
	"github.com/devbydaniel/release-notes-go/internal/ratelimit"
)

type loginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var loginRateLimiter = ratelimit.New(60, 5)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleLogin")
	userService := user.NewService(*user.NewRepository(h.DB))
	sessionService := session.NewService(*session.NewRepository(h.DB))

	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := loginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	h.log.Debug().Interface("req", req).Msg("Login request")

	// Check rate limit before proceeding
	if err := loginRateLimiter.Deduct(req.Email, 1); err != nil {
		h.log.Warn().Str("email", req.Email).Msg("Rate limit exceeded for login attempts")
		http.Error(w, "Too many login attempts. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// validate credentials
	user, err := userService.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, h.DB.ErrRecordNotFound) {
			http.Error(w, "Wrong credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error accessing user", http.StatusInternalServerError)
		return
	}
	if match := password.DoPasswordsMatch(user.Password, req.Password); !match {
		http.Error(w, "Wrong credentials", http.StatusUnauthorized)
		return
	}

	// create session
	token := sessionService.CreateToken()
	if err := sessionService.Create(token, user.ID); err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     session.AuthCookieName,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("HX-Redirect", "/release-notes")
	w.WriteHeader(http.StatusOK)
	return
}
