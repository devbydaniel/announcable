package passwordreset

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
)

// Handlers holds the dependencies for password reset handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the password reset page
type pageData struct {
	shared.BaseTemplateData
	Token string
}

var pageTmpl = templates.Construct(
	"invite-accept",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/reset-pw.html",
)

// ServeResetPasswordPage handles GET /reset-pw/{token}/
func (h *Handlers) ServeResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeResetPasswordPage")
	token := chi.URLParam(r, "token")

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Reset Password",
		},
		Token: token,
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
