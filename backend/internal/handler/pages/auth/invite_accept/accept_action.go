package invite_accept

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/internal/password"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
)

type acceptForm struct {
	Password string
	Confirm  string
}

// HandleAccept handles POST /invite-accept/{token}/
func (h *Handlers) HandleAccept(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleAccept")
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	token := chi.URLParam(r, "token")

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// decode form
	var acceptDTO acceptForm
	if err := schema.NewDecoder().Decode(&acceptDTO, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}
	h.deps.Log.Debug().Interface("acceptDTO", acceptDTO).Msg("acceptDTO")

	// validate form
	validate := validator.New()
	if err := validate.Struct(acceptDTO); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	invite, err := orgService.GetInviteWithToken(token)
	if err != nil {
		if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
			h.deps.Log.Debug().Msg("Invite not found")
			http.Error(w, "Invite not found", http.StatusNotFound)
			return
		}
		h.deps.Log.Error().Err(err).Msg("Error getting invite")
		http.Error(w, "Error getting invite", http.StatusInternalServerError)
		return
	}

	if err := password.IsValidPassword(acceptDTO.Password); err != nil {
		h.deps.Log.Debug().Err(err).Msg("Invalid password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if acceptDTO.Password != acceptDTO.Confirm {
		h.deps.Log.Debug().Msg("Passwords do not match")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	existingUser, err := userService.GetByEmail(invite.Email)
	if err != nil && !errors.Is(err, h.deps.DB.ErrRecordNotFound) {
		h.deps.Log.Error().Err(err).Msg("Error checking existing user")
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		h.deps.Log.Debug().Msg("User already exists")
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user, err := userService.Create(invite.Email, acceptDTO.Password, true)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	if err := orgService.AcceptInvite(invite, user); err != nil {
		userService.Delete(user.ID)
		http.Error(w, "Error accepting invite", http.StatusInternalServerError)
		return
	}

	successMsg := url.QueryEscape("invite accepted")
	w.Header().Set("HX-Redirect", "/login?success="+successMsg)
	w.WriteHeader(http.StatusCreated)
	return
}
