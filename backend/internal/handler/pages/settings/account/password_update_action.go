package account

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/user"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/internal/password"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// pwUpdateForm represents the password update form data
type pwUpdateForm struct {
	CurrentPassword string `schema:"current_password" validate:"required"`
	NewPassword     string `schema:"new_password" validate:"required"`
	Confirm         string `schema:"confirm" validate:"required"`
}

// HandlePasswordUpdate handles PATCH /settings/password
func (h *Handlers) HandlePasswordUpdate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandlePasswordUpdate")
	ctx := r.Context()
	userService := user.NewService(*user.NewRepository(h.deps.DB))

	userId, ok := ctx.Value(mw.UserIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Error updating password", http.StatusInternalServerError)
	}

	user, err := userService.GetById(uuid.MustParse(userId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error finding user")
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating password", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO pwUpdateForm
	if err := h.deps.Decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating password", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	if ok := password.DoPasswordsMatch(user.Password, updateDTO.CurrentPassword); !ok {
		h.deps.Log.Debug().Msg("Current password does not match")
		http.Error(w, "Current password does not match", http.StatusBadRequest)
		return
	}

	if err := password.IsValidPassword(updateDTO.NewPassword); err != nil {
		h.deps.Log.Debug().Str("pw", updateDTO.NewPassword).Msg("Password does not meet requirements")
		http.Error(w, "Password does not meet requirements", http.StatusBadRequest)
		return
	}

	if updateDTO.NewPassword != updateDTO.Confirm {
		h.deps.Log.Debug().Msg("Passwords do not match")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	if err := userService.UpdatePassword(uuid.MustParse(userId), updateDTO.NewPassword); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating password")
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
