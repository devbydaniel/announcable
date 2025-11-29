package register

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
)

// Handlers holds the dependencies for register handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the register page
type pageData struct {
	shared.BaseTemplateData
}

var pageTmpl = templates.Construct(
	"register",
	"layouts/root.html",
	"layouts/onboard.html",
	"pages/register.html",
)

// ServeRegisterPage handles GET /register/
func (h *Handlers) ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Register",
		},
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
