package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/session"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/devbydaniel/release-notes-go/templates"
)

type emailVerifyPageData struct {
	BaseTemplateData
}

var emailVerifyTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/verify-email.html",
)

func (h *Handler) HandleEmailVerifyPage(w http.ResponseWriter, r *http.Request) {
	sessionService := session.NewService(*session.NewRepository(h.DB))
	userService := user.NewService(*user.NewRepository(h.DB))
	token := r.URL.Query().Get("token")
	if token != "" {
		// user comes from the link in the email
		session, err := sessionService.ValidateSession(token)
		if err != nil {
			if errors.Is(err, h.DB.ErrRecordNotFound) {
				h.log.Warn().Str("token", token).Msg("Session not found")
				errorMsg := "Invalid token"
				escapedMsg := url.QueryEscape(errorMsg)
				redirectURL := fmt.Sprintf("/verify-email?error=%s", escapedMsg)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			}
			h.log.Error().Err(err).Msg("Error validating session")
			http.Error(w, "Error validating session", http.StatusInternalServerError)
			return
		}
		userId := session.UserID
		_, err = userService.GetById(userId)
		if err != nil {
			h.log.Error().Err(err).Msg("Error getting user")
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
		if err := userService.VerifyEmail(userId); err != nil {
			h.log.Error().Err(err).Msg("Error updating email verified")
			http.Error(w, "Error updating email verified", http.StatusInternalServerError)
			return
		}
		if err := sessionService.InvalidateUserSessions(userId); err != nil {
			h.log.Error().Err(err).Msg("Error invalidating user sessions")
		}
		successMsg := "Email verified"
		escapedMsg := url.QueryEscape(successMsg)
		redirectURL := fmt.Sprintf("/login?success=%s", escapedMsg)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	} else {
		// user comes from registration or other
		data := loginPageData{
			BaseTemplateData: BaseTemplateData{
				Title: "Verify Email",
			},
		}
		if err := emailVerifyTmpl.ExecuteTemplate(w, "root", data); err != nil {
			h.log.Error().Err(err).Msg("Error rendering page")
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
		return
	}
}
