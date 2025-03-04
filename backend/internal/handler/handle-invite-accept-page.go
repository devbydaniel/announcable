package handler

import (
	"errors"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
)

type InviteAcceptPage struct {
	BaseTemplateData
	Org   string
	Token string
	Email string
}

var inviteAcceptTmpl = templates.Construct(
	"invite-accept",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/invite-accept.html",
)

var inviteInvalidTmpl = templates.Construct(
	"invite-invalid",
	"layouts/root.html",
	"layouts/fullscreenmessage.html",
	"pages/invite-invalid.html",
)

func (h *Handler) HandleInviteAcceptPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleInviteAcceptPage")
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	token := chi.URLParam(r, "token")

	invite, err := organisationService.GetInviteWithToken(token)
	if err != nil {
		if errors.Is(err, h.DB.ErrRecordNotFound) {
			if err := inviteInvalidTmpl.ExecuteTemplate(w, "root", nil); err != nil {
				h.log.Error().Err(err).Msg("Error rendering page")
				http.Error(w, "Error rendering page", http.StatusInternalServerError)
				return
			}
			return
		}
		h.log.Error().Err(err).Msg("Error getting invite")
		http.Error(w, "Error getting invite", http.StatusInternalServerError)
		return
	}
	data := InviteAcceptPage{
		BaseTemplateData: BaseTemplateData{
			Title: "Login",
		},
		Org:   invite.Organisation.Name,
		Token: token,
		Email: invite.Email,
	}
	if err := inviteAcceptTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
