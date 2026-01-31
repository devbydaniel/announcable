package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var passwordResetTriggerRateLimiter = ratelimit.New(60, 5)

// HandlePasswordResetTrigger handles POST /users/{id}/password-reset
func (h *Handlers) HandlePasswordResetTrigger(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandlePasswordResetTrigger")
	ctx := r.Context()

	// Get admin's user ID for rate limiting
	adminUserID := ctx.Value(mw.UserIDKey).(string)
	if adminUserID == "" {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Error triggering password reset", http.StatusInternalServerError)
		return
	}

	// Check rate limit
	if err := passwordResetTriggerRateLimiter.Deduct(adminUserID, 1); err != nil {
		h.deps.Log.Warn().Str("admin_user_id", adminUserID).Msg("Rate limit exceeded for password reset trigger")
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// Get target OrgUser ID from URL
	orgUserID := chi.URLParam(r, "id")
	if orgUserID == "" {
		h.deps.Log.Error().Msg("OrgUser ID not found in URL")
		http.Error(w, "Error triggering password reset", http.StatusBadRequest)
		return
	}

	// Initialize services
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))

	// Get target OrgUser to find the actual user
	ou, err := orgService.GetOrgUser(uuid.MustParse(orgUserID))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting org user")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get target user
	targetUser, err := userService.GetByID(ou.UserID)
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error getting user")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Invalidate existing sessions
	if err := sessionService.InvalidateUserSessions(targetUser.ID); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error invalidating user sessions")
		// Continue - not critical
	}

	// Create password reset token (1 hour expiry)
	token := sessionService.CreateToken()
	if err := sessionService.CreateCustomDuration(token, targetUser.ID, 1*time.Hour); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error creating reset token")
		http.Error(w, "Error triggering password reset", http.StatusInternalServerError)
		return
	}

	// Build reset URL
	cfg := config.New()
	protocol := "https"
	if cfg.Env == "development" {
		protocol = "http"
	}
	resetURL := fmt.Sprintf("%s://%s/reset-pw/%s", protocol, cfg.BaseURL, token)

	if cfg.IsEmailEnabled() {
		// Send email
		if err := userService.SendPwResetEmail(targetUser, token); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error sending password reset email")
			http.Error(w, "Error sending password reset email", http.StatusInternalServerError)
			return
		}
		// Redirect with success message
		successMsg := "password reset email sent"
		escapedMsg := url.QueryEscape(successMsg)
		redirectURL := fmt.Sprintf("/users?success=%s", escapedMsg)
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
	} else {
		// Return reset URL as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"resetURL": resetURL})
	}
}
