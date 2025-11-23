package detail

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
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
	CtaUrlOverrideIsChecked      bool
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
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	hasActiveSubscription, ok := r.Context().Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.deps.Log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	h.deps.Log.Debug().Str("id", id).Str("orgId", orgId).Msg("Release note page")

	// get release note
	rn, err := releaseNoteService.GetOne(id, orgId)
	if err != nil {
		http.Error(w, "Error getting release note", http.StatusInternalServerError)
	}

	cfg := config.New()

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title:                 rn.Title,
			HasActiveSubscription: hasActiveSubscription,
			ShowSubscriptionUI:    cfg.IsCloud(),
		},
		Rn:                           rn,
		IsEdit:                       true,
		TextWebsiteOverrideIsChecked: rn.DescriptionLong != "",
		HideCtaIsChecked:             rn.HideCta,
		CtaLabelOverrideIsChecked:    rn.CtaLabelOverride != "",
		CtaUrlOverrideIsChecked:      rn.CtaUrlOverride != "",
	}
	h.deps.Log.Debug().Interface("data", data).Msg("Data")
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
