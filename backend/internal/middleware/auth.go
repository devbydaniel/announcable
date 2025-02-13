package mw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	"github.com/devbydaniel/release-notes-go/internal/domain/session"
)

type contextKey string

const (
	SessionIdKey     contextKey = "sessionId"
	UserIDKey        contextKey = "userId"
	OrgRoleKey       contextKey = "orgRole"
	OrgIDKey         contextKey = "orgId"
	EmailVerifiedKey contextKey = "emailVerified"
)

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	sessionService := session.NewService(*session.NewRepository(h.DB))
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw Authenticate")
		// get session cookie
		cookie, err := r.Cookie("releasebeacon-session")
		if err != nil {
			escapedMsg := url.QueryEscape("please log in")
			url := fmt.Sprintf("/login?info=%s", escapedMsg)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}

		authToken := cookie.Value

		session, err := sessionService.ValidateSession(authToken)
		if err != nil {
			if errors.Is(err, h.DB.ErrRecordNotFound) {
				escapedMsg := url.QueryEscape("please log in")
				url := fmt.Sprintf("/login?info=%s", escapedMsg)
				http.Redirect(w, r, url, http.StatusSeeOther)
				return
			}
			http.Error(w, "Error validating session", http.StatusInternalServerError)
			return
		}

		ou, err := orgService.GetOrgUserByUserId(session.UserID)
		if err != nil {
			http.Error(w, "Error getting organisation user", http.StatusInternalServerError)
			return
		}
		if ou == nil {
			sessionService.InvalidateUserSessions(session.UserID)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.log.Debug().Interface("ou", ou).Msg("OrganisationUser")

		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionIdKey, session.ID.String())
		ctx = context.WithValue(ctx, EmailVerifiedKey, ou.User.EmailVerified)
		ctx = context.WithValue(ctx, UserIDKey, session.UserID.String())
		ctx = context.WithValue(ctx, OrgRoleKey, ou.Role)
		ctx = context.WithValue(ctx, OrgIDKey, ou.OrganisationID.String())
		r = r.WithContext(ctx)
		h.log.Trace().Str("userId", session.UserID.String()).Str("role", ou.Role.String()).Str("orgId", ou.OrganisationID.String()).Msg("Authenticated")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Authorize(permissions ...rbac.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.log.Trace().Msg("mw Authorize")
			ctx := r.Context()
			emailVerified, ok := ctx.Value(EmailVerifiedKey).(bool)
			if !ok {
				h.log.Warn().Msg("EmailVerified not found in context")
				http.Error(w, "EmailVerified not found in context", http.StatusInternalServerError)
				return
			}
			if !emailVerified {
				h.log.Warn().Msg("Email not verified")
				http.Redirect(w, r, "/verify-email", http.StatusSeeOther)
				return
			}
			orgRole, ok := ctx.Value(OrgRoleKey).(rbac.Role)
			if !ok {
				h.log.Warn().Msg("Role not found in context")
				http.Error(w, "Role not found in context", http.StatusInternalServerError)
				return
			}
			for _, permission := range permissions {
				if !rbac.HasPermission(orgRole, permission) {
					h.log.Warn().Str("role", orgRole.String()).Str("permission", permission.String()).Msg("Unauthorized")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			h.log.Trace().Str("role", orgRole.String()).Msg("Authorized")
			next.ServeHTTP(w, r)
		})
	}
}
