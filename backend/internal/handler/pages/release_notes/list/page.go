package list

import (
	"net/http"
	"strconv"
	"time"

	releasenotelikes "github.com/devbydaniel/announcable/internal/domain/release-note-likes"
	releasenotemetrics "github.com/devbydaniel/announcable/internal/domain/release-note-metrics"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/handler/shared"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/devbydaniel/announcable/templates"
)

// Handlers holds the dependencies for release notes list handlers
type Handlers struct {
	deps *shared.Dependencies
}

// New creates a new Handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{deps: deps}
}

// ReleaseNoteWithMetrics represents a release note with its associated metrics
type ReleaseNoteWithMetrics struct {
	*releasenotes.ReleaseNote
	ViewCount     int
	LikeCount     int
	CtaClickCount int
}

// pageData holds the template data for the release notes list page
type pageData struct {
	shared.BaseTemplateData
	ReleaseNotes []*ReleaseNoteWithMetrics
	NextPageLink string
	PrevPageLink string
}

var pageTmpl = templates.Construct(
	"release-notes",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/release-notes-list.html",
)

// ServeReleaseNotesListPage handles GET /release-notes/
func (h *Handlers) ServeReleaseNotesListPage(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("ServeReleaseNotesListPage")
	ctx := r.Context()
	orgID, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.deps.DB, h.deps.ObjStore))
	metricsService := releasenotemetrics.NewService(releasenotemetrics.NewRepository(h.deps.DB))
	likesService := releasenotelikes.NewService(releasenotelikes.NewRepository(h.deps.DB))

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
		h.deps.Log.Error().Err(err).Msg("Error parsing page")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing pageSize")
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}

	releaseNotes, err := releaseNotesService.GetAll(orgID, pageInt, pageSizeInt)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// Create slice to hold release notes with metrics
	releaseNotesWithMetrics := make([]*ReleaseNoteWithMetrics, len(releaseNotes.Items))

	// For each release note, get its metrics and format date
	for i, rn := range releaseNotes.Items {
		// Format date as before
		if rn.ReleaseDate != nil {
			releaseDate, err := time.Parse("2006-01-02", *rn.ReleaseDate)
			if err != nil {
				h.deps.Log.Error().Err(err).Msg("Error parsing release date")
				continue
			}
			rd := releaseDate.Format("02.01.2006")
			rn.ReleaseDate = &rd
		} else {
			rd := ""
			rn.ReleaseDate = &rd
		}

		// Get view count for this release note
		viewCount, err := metricsService.GetViewCount(rn.ID)
		if err != nil {
			h.deps.Log.Error().Err(err).Msg("Error getting view count for release note")
			viewCount = 0
		}

		// Get like count for this release note
		likeCount, err := likesService.GetLikeCount(rn.ID)
		if err != nil {
			h.deps.Log.Error().Err(err).Msg("Error getting like count for release note")
			likeCount = 0
		}

		// Get CTA click count for this release note
		ctaClickCount, err := metricsService.GetCtaClickCount(rn.ID)
		if err != nil {
			h.deps.Log.Error().Err(err).Msg("Error getting CTA click count for release note")
			ctaClickCount = 0
		}

		releaseNotesWithMetrics[i] = &ReleaseNoteWithMetrics{
			ReleaseNote:   rn,
			ViewCount:     viewCount,
			LikeCount:     likeCount,
			CtaClickCount: ctaClickCount,
		}
	}

	nextPageLink := ""
	if releaseNotes.Page < releaseNotes.TotalPages {
		nextPageLink = "/release-notes?page=" + strconv.Itoa(releaseNotes.Page+1) + "&pageSize=" + pageSize
	}
	prevPageLink := ""
	if releaseNotes.Page > 1 {
		prevPageLink = "/release-notes?page=" + strconv.Itoa(releaseNotes.Page-1) + "&pageSize=" + pageSize
	}

	data := pageData{
		BaseTemplateData: shared.BaseTemplateData{
			Title: "Release Notes",
		},
		ReleaseNotes: releaseNotesWithMetrics,
		NextPageLink: nextPageLink,
		PrevPageLink: prevPageLink,
	}
	if err := pageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
