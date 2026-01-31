package detail

import (
	"net/http"

	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var publishButtonTmpl = templates.Construct("publish-btn", "partials/hx-publish-rn-button.html")

// HandleReleaseNotePublish handles PATCH /release-notes/{id}/publish
func (h *Handlers) HandleReleaseNotePublish(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleReleaseNotePublish")
	ctx := r.Context()
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.deps.DB, h.deps.ObjStore))

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing ID")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	userID := ctx.Value(mw.UserIDKey).(string)
	if userID == "" {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	orgID := ctx.Value(mw.OrgIDKey).(string)
	if orgID == "" {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	shouldPublish := r.FormValue("publish") == "true"

	h.deps.Log.Debug().Interface("publishDTO", shouldPublish).Msg("publishDTO")
	if err := releaseNotesService.ChangePublishedStatus(id, shouldPublish); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating release note")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
		return
	}

	var templateName string
	if shouldPublish {
		templateName = "hx-unpublish-rn-button"
	} else {
		templateName = "hx-publish-rn-button"
	}

	err = publishButtonTmpl.ExecuteTemplate(w, templateName, id.String())
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering template")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
		return
	}
}
