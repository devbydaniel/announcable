package handler

import (
	"net/http"
	"time"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

type releaseNotesPageData struct {
	Title        string
	ReleaseNotes []*releasenotes.ReleaseNote
}

var releaseNotesPageTmpl = templates.Construct(
	"release-notes",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-notes-list.html",
)

func (h *Handler) HandleReleaseNotesPage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNotesPage")
	ctx := r.Context()
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
	releaseNotes, err := releaseNotesService.GetAll(orgId)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// format release date
	for _, rn := range releaseNotes {
		if rn.ReleaseDate != nil {
			releaseDate, err := time.Parse("2006-01-02", *rn.ReleaseDate)
			if err != nil {
				h.log.Error().Err(err).Msg("Error parsing release date")
				continue
			}
			rd := releaseDate.Format("02.01.2006")
			rn.ReleaseDate = &rd
		} else {
			rd := ""
			rn.ReleaseDate = &rd
		}
	}
	data := releaseNotesPageData{
		Title:        "Release Notes",
		ReleaseNotes: releaseNotes,
	}
	if err := releaseNotesPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
