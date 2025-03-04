package handler

import (
	"net/http"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

var releaseNoteCreatePageTmpl = templates.Construct(
	"new-release-note",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-note-create-edit.html",
)

type releaseNotePageData struct {
	BaseTemplateData
	Rn                           *releasenotes.ReleaseNote
	IsEdit                       bool
	TextWebsiteOverrideIsChecked bool
	HideCtaIsChecked             bool
	CtaLabelOverrideIsChecked    bool
	CtaUrlOverrideIsChecked      bool
}

func (h *Handler) HandleReleaseNoteCreatePage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNoteCreatePage")
	hasActiveSubscription, ok := r.Context().Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	data := releaseNotePageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "New Release Note",
			HasActiveSubscription: hasActiveSubscription,
		},
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
