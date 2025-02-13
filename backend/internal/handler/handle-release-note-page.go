package handler

import (
	"net/http"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
)

var releaseNoteDetailPageTmpl = templates.Construct(
	"new-release-note",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-note-create-edit.html",
)

func (h *Handler) HandleReleaseNotePage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNotePage")
	releaseNoteService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))

	id := chi.URLParam(r, "id")
	h.log.Debug().Str("id", id).Msg("id URL param")
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	h.log.Debug().Str("id", id).Str("orgId", orgId).Msg("Release note page")

	// get release note
	rn, err := releaseNoteService.GetOne(id, orgId)
	if err != nil {
		http.Error(w, "Error getting release note", http.StatusInternalServerError)
	}

	data := releaseNotePageData{
		Title:                        rn.Title,
		Rn:                           rn,
		IsEdit:                       true,
		TextWebsiteOverrideIsChecked: rn.DescriptionLong != "",
		HideCtaIsChecked:             rn.HideCta,
		CtaLabelOverrideIsChecked:    rn.CtaLabelOverride != "",
		CtaUrlOverrideIsChecked:      rn.CtaUrlOverride != "",
	}
	h.log.Debug().Interface("data", data).Msg("Data")
	if err := releaseNoteDetailPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
