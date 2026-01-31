package detail

import (
	"net/http"

	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
)

// Handlers holds the dependencies for release note detail handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// pageData holds the template data for the release note detail/edit page
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

// ServeReleaseNoteDetailPage handles GET /release-notes/{id}
func (h *Handlers) ServeReleaseNoteDetailPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeReleaseNoteDetailPage")
	releaseNoteService := releasenotes.NewService(*releasenotes.NewRepository(h.deps.DB, h.deps.ObjStore))

	id := chi.URLParam(r, "id")
	h.deps.Log.Debug().Str("id", id).Msg("id URL param")
	orgID, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	h.deps.Log.Debug().Str("id", id).Str("orgID", orgID).Msg("Release note page")

	// get release note
	rn, err := releaseNoteService.GetOne(id, orgID)
	if err != nil {
		http.Error(w, "Error getting release note", http.StatusInternalServerError)
	}

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: rn.Title,
		},
		Rn:                           rn,
		IsEdit:                       true,
		TextWebsiteOverrideIsChecked: rn.DescriptionLong != "",
		HideCtaIsChecked:             rn.HideCta,
		CtaLabelOverrideIsChecked:    rn.CtaLabelOverride != "",
		CtaURLOverrideIsChecked:      rn.CtaURLOverride != "",
	}
	h.deps.Log.Debug().Interface("data", data).Msg("Data")
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
