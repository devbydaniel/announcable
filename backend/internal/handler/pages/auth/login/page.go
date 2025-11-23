package login

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
	"github.com/devbydaniel/release-notes-go/templates"
)

// Handlers holds the dependencies for login handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the login page
type pageData struct {
	shared.BaseTemplateData
}

var pageTmpl = templates.Construct(
	"login",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/login.html",
)

// ServeLoginPage handles GET /login/
func (h *Handlers) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Login",
		},
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
