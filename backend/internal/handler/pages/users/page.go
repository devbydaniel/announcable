package users

import (
	"net/http"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

// Handlers holds the dependencies for users list handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// UserData represents user information for the page
type UserData struct {
	UserID    string
	OrgUserID string
	Email     string
	Role      string
}

// InviteData represents invite information for the page
type InviteData struct {
	ID        string
	Email     string
	Role      string
	IsExpired bool
}

// pageData holds the template data for the users list page
type pageData struct {
	shared.BaseTemplateData
	Users   []*UserData
	Invites []*InviteData
	OwnID   string
}

var pageTmpl = templates.Construct(
	"users",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/user-list.html",
)

// ServeUsersPage handles GET /users/
func (h *Handlers) ServeUsersPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeUsersPage")
	ctx := r.Context()
	userID, ok := ctx.Value(mw.UserIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	orgService := organisation.NewService(*organisation.NewRepository(h.deps.DB))

	orgUsers, err := orgService.GetOrgUsers(uuid.MustParse(orgId))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	userData := make([]*UserData, 0)
	for _, ou := range orgUsers {
		userData = append(userData, &UserData{
			OrgUserID: ou.ID.String(),
			UserID:    ou.User.ID.String(),
			Email:     ou.User.Email,
			Role:      ou.Role.String(),
		})
	}

	invites, err := orgService.GetInvites(uuid.MustParse(orgId))
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	inviteData := make([]*InviteData, 0)
	for _, i := range invites {
		inviteData = append(inviteData, &InviteData{
			ID:        i.ID.String(),
			Email:     i.Email,
			Role:      i.Role.String(),
			IsExpired: time.Now().After(time.Unix(i.ExpiresAt, 0)),
		})
	}

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Users",
		},
		Users:   userData,
		Invites: inviteData,
		OwnID:   userID,
	}
	h.deps.Log.Debug().Interface("pageData", data).Msg("Page data")
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
