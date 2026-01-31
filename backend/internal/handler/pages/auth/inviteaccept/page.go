package inviteaccept

import (
	"errors"
	"net/http"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
)

// Handlers holds the dependencies for invite accept handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the invite accept page
type pageData struct {
	shared.BaseTemplateData
	Org   string
	Token string
	Email string
}

var pageTmpl = templates.Construct(
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

// ServeInviteAcceptPage handles GET /invite-accept/{token}/
func (h *Handlers) ServeInviteAcceptPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeInviteAcceptPage")
	organisationService := organisation.NewService(*organisation.NewRepository(h.deps.DB))
	token := chi.URLParam(r, "token")

	invite, err := organisationService.GetInviteWithToken(token)
	if err != nil {
		if errors.Is(err, h.deps.DB.ErrRecordNotFound) {
			if err := inviteInvalidTmpl.ExecuteTemplate(w, "root", nil); err != nil {
				h.deps.Log.Error().Err(err).Msg("Error rendering page")
				http.Error(w, "Error rendering page", http.StatusInternalServerError)
				return
			}
			return
		}
		h.deps.Log.Error().Err(err).Msg("Error getting invite")
		http.Error(w, "Error getting invite", http.StatusInternalServerError)
		return
	}
	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Login",
		},
		Org:   invite.Organisation.Name,
		Token: token,
		Email: invite.Email,
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
