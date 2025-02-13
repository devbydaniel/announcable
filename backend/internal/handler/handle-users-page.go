package handler

import (
	"net/http"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/google/uuid"
)

type UserData struct {
	UserID    string
	OrgUserID string
	Email     string
	Role      string
}

type InviteData struct {
	ID        string
	Email     string
	Role      string
	IsExpired bool
}

type usersPageData struct {
	Title   string
	Users   []*UserData
	Invites []*InviteData
	OwnID   string
}

var usersPageTmpl = templates.Construct(
	"users",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/user-list.html",
)

func (h *Handler) HandleUsersPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleUsersPage")
	ctx := r.Context()
	userID, ok := ctx.Value(mw.UserIDKey).(string)
	if !ok {
		h.log.Error().Msg("User ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))

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

	pageData := usersPageData{
		Title:   "Users",
		Users:   userData,
		Invites: inviteData,
		OwnID:   userID,
	}
	h.log.Debug().Interface("pageData", pageData).Msg("Page data")
	if err := usersPageTmpl.ExecuteTemplate(w, "root", pageData); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
