package verifyemail

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/announcable/internal/domain/session"
	"github.com/devbydaniel/announcable/internal/domain/user"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
)

// Handlers holds the dependencies for verify email handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the verify email page
type pageData struct {
	shared.BaseTemplateData
}

var pageTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/verify-email.html",
)

// ServeVerifyEmailPage handles GET /verify-email/
func (h *Handlers) ServeVerifyEmailPage(w http.ResponseWriter, r *http.Request) {
	sessionService := session.NewService(*session.NewRepository(h.deps.DB))
	userService := user.NewService(*user.NewRepository(h.deps.DB))
	token := r.URL.Query().Get("token")
	if token != "" {
		// user comes from the link in the email
		session, err := sessionService.ValidateSession(token)
		if err != nil {
			if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
				h.deps.Log.Warn().Str("token", token).Msg("Session not found")
				errorMsg := "Invalid token"
				escapedMsg := url.QueryEscape(errorMsg)
				redirectURL := fmt.Sprintf("/verify-email?error=%s", escapedMsg)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			}
			h.deps.Log.Error().Err(err).Msg("Error validating session")
			http.Error(w, "Error validating session", http.StatusInternalServerError)
			return
		}
		userID := session.UserID
		_, err = userService.GetByID(userID)
		if err != nil {
			h.deps.Log.Error().Err(err).Msg("Error getting user")
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
		if err := userService.VerifyEmail(userID); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error updating email verified")
			http.Error(w, "Error updating email verified", http.StatusInternalServerError)
			return
		}
		if err := sessionService.InvalidateUserSessions(userID); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error invalidating user sessions")
		}
		successMsg := "Email verified"
		escapedMsg := url.QueryEscape(successMsg)
		redirectURL := fmt.Sprintf("/login?success=%s", escapedMsg)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	} else {
		// user comes from registration or other
		data := pageData{
			BaseTemplateData: shared.BaseTemplateData{
				Title: "Verify Email",
			},
		}
		if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
			h.deps.Log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
		return
	}
}
