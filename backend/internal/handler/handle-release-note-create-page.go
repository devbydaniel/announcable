package handler

import (
	"net/http"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/devbydaniel/release-notes-go/templates"
)

var releaseNoteCreatePageTmpl = templates.Construct(
	"new-release-note",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-note-create-edit.html",
)

type releaseNotePageData struct {
	Title                        string
	Rn                           *releasenotes.ReleaseNote
	IsEdit                       bool
	TextWebsiteOverrideIsChecked bool
	HideCtaIsChecked             bool
	CtaLabelOverrideIsChecked    bool
	CtaUrlOverrideIsChecked      bool
}

func (h *Handler) HandleReleaseNoteCreatePage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNoteCreatePage")
	data := releaseNotePageData{
		Title:                        "New Release Note",
		Rn:                           &releasenotes.ReleaseNote{},
		IsEdit:                       false,
		TextWebsiteOverrideIsChecked: false,
		HideCtaIsChecked:             false,
		CtaLabelOverrideIsChecked:    false,
		CtaUrlOverrideIsChecked:      false,
	}

	if err := releaseNoteCreatePageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
