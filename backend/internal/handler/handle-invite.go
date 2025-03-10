package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/ratelimit"
	"github.com/go-playground/validator"
)

type userInviteForm struct {
	Email string    `json:"email" validate:"required,email"`
	Role  rbac.Role `json:"role" validate:"required"`
}

var inviteUserRateLimiter = ratelimit.New(60, 10)

func (h *Handler) HandleInvite(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleInvite")
	ctx := r.Context()

	userId := ctx.Value(mw.UserIDKey).(string)
	if userId == "" {
		h.log.Error().Msg("User ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	// check rate limit
	if err := inviteUserRateLimiter.Deduct(userId, 1); err != nil {
		h.log.Warn().Str("user_id", userId).Msg("Rate limit exceeded for invite user requests")
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
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

	successMsg := "invite sent"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := fmt.Sprintf("/users?success=%s", escapedMsg)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
}
