package handler

import (
	"net/http"
	"strconv"
	"time"

	releasenotemetrics "github.com/devbydaniel/release-notes-go/internal/domain/release-note-metrics"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

type releaseNoteWithMetrics struct {
	*releasenotes.ReleaseNote
	ViewCount int
}

type releaseNotesPageData struct {
	BaseTemplateData
	ReleaseNotes []*releaseNoteWithMetrics
	NextPageLink string
	PrevPageLink string
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
	hasActiveSubscription, ok := ctx.Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
	metricsService := releasenotemetrics.NewService(releasenotemetrics.NewRepository(h.DB))
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

	releaseNotes, err := releaseNotesService.GetAll(orgId, pageInt, pageSizeInt)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Create slice to hold release notes with metrics
	releaseNotesWithMetrics := make([]*releaseNoteWithMetrics, len(releaseNotes.Items))

	// For each release note, get its metrics and format date
	for i, rn := range releaseNotes.Items {
		// Format date as before
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

		// Get view count for this release note
		viewCount, err := metricsService.GetViewCount(rn.ID)
		if err != nil {
			h.log.Error().Err(err).Msg("Error getting view count for release note")
			viewCount = 0
		}

		releaseNotesWithMetrics[i] = &releaseNoteWithMetrics{
			ReleaseNote: rn,
			ViewCount:   viewCount,
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
	data := releaseNotesPageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Release Notes",
			HasActiveSubscription: hasActiveSubscription,
		},
		ReleaseNotes: releaseNotesWithMetrics,
		NextPageLink: nextPageLink,
		PrevPageLink: prevPageLink,
	}
	if err := releaseNotesPageTmpl.ExecuteTemplate(w, "root", data); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
