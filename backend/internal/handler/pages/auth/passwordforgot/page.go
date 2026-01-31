package passwordforgot

import (
	"net/http"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
)

// Handlers holds the dependencies for password forgot handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the password forgot page
type pageData struct {
	shared.BaseTemplateData
	EmailEnabled bool
}

var pageTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/forgot-pw.html",
)

// ServeForgotPasswordPage handles GET /forgot-pw/
func (h *Handlers) ServeForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeForgotPasswordPage")
	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Password reset",
		},
		EmailEnabled: config.New().IsEmailEnabled(),
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
