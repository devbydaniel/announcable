package users

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/ratelimit"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// userInviteForm represents the invite form data
type userInviteForm struct {
	Email string    `json:"email" validate:"required,email"`
	Role  rbac.Role `json:"role" validate:"required"`
}

var inviteUserRateLimiter = ratelimit.New(60, 10)

// HandleInviteCreate handles POST /invites/
func (h *Handlers) HandleInviteCreate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleInviteCreate")
	ctx := r.Context()

	userId := ctx.Value(mw.UserIDKey).(string)
	if userId == "" {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Error creating invite", http.StatusInternalServerError)
	}

	// check rate limit
	if err := inviteUserRateLimiter.Deduct(userId, 1); err != nil {
		h.deps.Log.Warn().Str("user_id", userId).Msg("Rate limit exceeded for invite user requests")
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error creating invite", http.StatusInternalServerError)
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error creating invite", http.StatusBadRequest)
		return
	}

	// decode form
	var inviteDTO userInviteForm
	if err := h.deps.Decoder.Decode(&inviteDTO, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error creating invite", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(inviteDTO); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error creating invite", http.StatusBadRequest)
		return
	}

	// create invite and send email
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	if _, err := orgService.InviteUser(uuid.MustParse(orgId), inviteDTO.Email, inviteDTO.Role); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error creating invite")
		http.Error(w, "Error creating invite", http.StatusInternalServerError)
		return
	}

	successMsg := "invite sent"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := fmt.Sprintf("/users?success=%s", escapedMsg)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
}
