package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	"github.com/devbydaniel/release-notes-go/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type releaseNotesWebsiteData struct {
	Cfg *releasepageconfig.ReleasePageConfig
	Rns []*releasenotes.ReleaseNote
}

var releaseNotesWebsiteTmpl = templates.Construct("release-notes-website", "pages/release-notes-website.html")

func (h *Handler) HandleReleasePage(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleasePage")
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	lpconfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))
	rnService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))

	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.log.Error().Msg("Org ID not found in URL")
		http.Error(w, "Error getting widget config", http.StatusBadRequest)
		return
	}

	org, err := organisationService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting org ID")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	config, err := lpconfigService.Get(org.ID.String())
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	rns, err := rnService.GetAllWithImgUrl(org.ID.String())
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release notes")
		http.Error(w, "Error getting release notes", http.StatusInternalServerError)
		return
	}

	data := releaseNotesWebsiteData{
		Cfg: config,
		Rns: rns,
	}

	if err := releaseNotesWebsiteTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
