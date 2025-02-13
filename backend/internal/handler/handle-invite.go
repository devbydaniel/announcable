package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type userInviteForm struct {
	Email string    `json:"email" validate:"required,email"`
	Role  rbac.Role `json:"role" validate:"required"`
}

func (h *Handler) HandleInvite(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleInvite")
	ctx := r.Context()
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))

	userId := ctx.Value(mw.UserIDKey).(string)
	if userId == "" {
		h.log.Error().Msg("User ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// decode form
	var inviteDTO userInviteForm
	if err := h.decoder.Decode(&inviteDTO, r.PostForm); err != nil {
		h.log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(inviteDTO); err != nil {
		h.log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// create invite
	orgService.InviteUser(uuid.MustParse(orgId), inviteDTO.Email, rbac.Role(inviteDTO.Role))

	successMsg := "invite sent"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := fmt.Sprintf("/users?success=%s", escapedMsg)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
}
