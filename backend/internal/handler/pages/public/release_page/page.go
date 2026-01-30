package release_page

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	"github.com/devbydaniel/announcable/internal/util"
	"github.com/devbydaniel/announcable/templates"
	"github.com/go-chi/chi/v5"
)

// Handlers provides public release page handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new public release page handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}

// ReleaseNotesWebsiteData represents the public release page template data
type ReleaseNotesWebsiteData struct {
	Cfg *releasepageconfig.ReleasePageConfig
	Rns []*releasenotes.ReleaseNote
}

var releaseNotesWebsiteTmpl = templates.Construct("release-notes-website", "pages/release-notes-website.html")

// ServeReleasePage renders the public release notes page
func (h *Handlers) ServeReleasePage(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("ServeReleasePage")
	backLinkLabel := r.URL.Query().Get("backLinkLabel")
	backLinkUrl := r.URL.Query().Get("backLinkUrl")
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
		h.Log.Error().Err(err).Msg("Error parsing page")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error parsing pageSize")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}

	orgSlug := chi.URLParam(r, "orgSlug")
	if orgSlug == "" {
		h.Log.Error().Msg("Org slug not found in URL")
		http.Error(w, "Error getting widget config", http.StatusBadRequest)
		return
	}

	config, err := releasePageConfigService.GetBySlug(orgSlug)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	if config.DisableReleasePage {
		http.NotFound(w, r)
		return
	}

	org, err := organisationService.GetOrg(config.OrganisationID)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting org ID")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	filters := map[string]interface{}{
		"is_published":         true,
		"hide_on_release_page": false,
	}
	rns, err := rnService.GetAllWithImgUrl(org.ID.String(), pageInt, pageSizeInt, filters)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting release notes")
		http.Error(w, "Error getting release notes", http.StatusInternalServerError)
		return
	}

	// format release date and transform media links
	for _, rn := range rns.Items {
		if rn.ReleaseDate != nil {
			releaseDate, err := time.Parse("2006-01-02", *rn.ReleaseDate)
			if err != nil {
				h.Log.Error().Err(err).Msg("Error parsing release date")
				continue
			}
			rd := releaseDate.Format("02.01.2006")
			rn.ReleaseDate = &rd
		} else {
			rd := ""
			rn.ReleaseDate = &rd
		}

		// Transform media link to embed URL
		if rn.MediaLink != "" {
			rn.MediaLink = util.TransformMediaLink(rn.MediaLink)
		}
	}

	// adjust back link if there's query params
	if backLinkUrl != "" {
		config.BackLinkUrl = url.QueryEscape(backLinkUrl)
	}
	if backLinkLabel != "" {
		config.BackLinkLabel = url.QueryEscape(backLinkLabel)
	}

	data := ReleaseNotesWebsiteData{
		Cfg: config,
		Rns: rns.Items,
	}

	if err := releaseNotesWebsiteTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
