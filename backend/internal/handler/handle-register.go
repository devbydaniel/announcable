package handler

import (
	"errors"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/internal/password"
)

type registerForm struct {
	OrgName         string
	Email           string
	Password        string
	ConfirmPassword string
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleRegister")
	userService := user.NewService(*user.NewRepository(h.DB))
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))
	sessionService := session.NewService(*session.NewRepository(h.DB))

	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := registerForm{
		OrgName:         r.FormValue("org"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm"),
	}
	h.log.Debug().Interface("req", req).Msg("Register request")

	if err := password.IsValidPassword(req.Password); err != nil {
		h.log.Debug().Err(err).Msg("Invalid password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		h.log.Debug().Msg("Passwords do not match")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	existing, err := userService.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, h.DB.ErrRecordNotFound) {
		h.log.Debug().Err(err).Msg("Error creating user")
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}
	if existing != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	user, err := userService.Create(req.Email, req.Password, false)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	_, err = orgService.CreateOrgWithAdmin(req.OrgName, user)
	if err != nil {
		userService.Delete(user.ID)
		http.Error(w, "Error creating organisation", http.StatusInternalServerError)
		return
	}

	token := sessionService.CreateToken()
	if err := sessionService.Create(token, user.ID); err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	if err := userService.SendVerifcationEmail(user, token); err != nil {
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/verify-email")
	w.WriteHeader(http.StatusCreated)
}
