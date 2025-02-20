package handler

import (
	"net/http"
	"strconv"

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
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))
	rnService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageSize := r.URL.Query().Get("pageSize")
	if pageSize == "" {
		pageSize = "10"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		h.log.Error().Err(err).Msg("Error parsing page")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		h.log.Error().Err(err).Msg("Error parsing pageSize")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}

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

	config, err := releasePageConfigService.Get(org.ID)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	rns, err := rnService.GetAllWithImgUrl(org.ID.String(), pageInt, pageSizeInt)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release notes")
		http.Error(w, "Error getting release notes", http.StatusInternalServerError)
		return
	}

	data := releaseNotesWebsiteData{
		Cfg: config,
		Rns: rns.Items,
	}

	if err := releaseNotesWebsiteTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
