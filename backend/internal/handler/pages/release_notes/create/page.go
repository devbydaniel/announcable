package create

import (
	"net/http"

	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/templates"
)

// Handlers holds the dependencies for release note create handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the release note create/edit page
type pageData struct {
	shared.BaseTemplateData
	Rn                           *releasenotes.ReleaseNote
	IsEdit                       bool
	TextWebsiteOverrideIsChecked bool
	HideCtaIsChecked             bool
	CtaLabelOverrideIsChecked    bool
	CtaURLOverrideIsChecked      bool
}

var pageTmpl = templates.Construct(
	"new-release-note",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-note-create-edit.html",
)

// ServeReleaseNoteCreatePage handles GET /release-notes/new
func (h *Handlers) ServeReleaseNoteCreatePage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeReleaseNoteCreatePage")

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "New Release Note",
		},
		Rn:                           &releasenotes.ReleaseNote{},
		IsEdit:                       false,
		TextWebsiteOverrideIsChecked: false,
		HideCtaIsChecked:             false,
		CtaLabelOverrideIsChecked:    false,
		CtaURLOverrideIsChecked:      false,
	}

	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
