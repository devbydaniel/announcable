package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/devbydaniel/release-notes-go/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type serveReleaseNotesWidgetResponseBodyReleaseNotes struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Date               string `json:"date"`
	ImageSrc           string `json:"imageSrc"`
	MediaLink          string `json:"media_link"`
	Text               string `json:"text"`
	LastUpdateOn       string `json:"last_update_on"`
	CtaLabelOverride   string `json:"cta_label_override"`
	CtaHrefOverride    string `json:"cta_href_override"`
	HideCta            bool   `json:"hide_cta"`
	AttentionMechanism string `json:"attentionMechanism"`
}

type serveReleaseNotesWidgetResponseBody struct {
	Data []serveReleaseNotesWidgetResponseBodyReleaseNotes `json:"data"`
}

func (h *Handler) HandleReleaseNotesServe(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNotesServe")
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
	forWidgetOrWebsite := r.URL.Query().Get("for")
	h.log.Debug().Str("for", forWidgetOrWebsite).Msg("For widget or website")
	h.log.Debug().Msg("Getting page and pageSize")
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
		http.Error(w, "Error getting release notes", http.StatusBadRequest)
		return
	}

	org, err := organisationService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting org ID")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}
	filters := map[string]interface{}{
		"is_published": true,
	}
	if forWidgetOrWebsite == "widget" {
		filters["hide_on_widget"] = false
	}
	if forWidgetOrWebsite == "website" {
		filters["hide_on_release_page"] = false
	}

	releaseNotes, err := releaseNotesService.GetAllWithImgUrl(org.ID.String(), pageInt, pageSizeInt, filters)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release notes")
		http.Error(w, "Error getting release notes", http.StatusInternalServerError)
		return
	}
	h.log.Debug().Int("releaseNotes", len(releaseNotes.Items)).Msg("Number of release notes")

	var res serveReleaseNotesWidgetResponseBody
	for _, rn := range releaseNotes.Items {
		if rn.IsPublished == false {
			continue
		}
		var releaseDate string
		if rn.ReleaseDate != nil {
			parsedDate, err := time.Parse("2006-01-02", *rn.ReleaseDate)
			if err != nil {
				h.log.Warn().Err(err).Msg("Error parsing date")
			} else {
				releaseDate = parsedDate.Format("02.01.2006")
			}
		} else {
			releaseDate = ""
		}
		h.log.Debug().Str("releaseDate", releaseDate).Msg("Release date")
		if rn.MediaLink != "" {
			rn.MediaLink = util.TransformMediaLink(rn.MediaLink)
		}
		res.Data = append(res.Data, serveReleaseNotesWidgetResponseBodyReleaseNotes{
			ID:                 rn.ID.String(),
			Title:              rn.Title,
			Date:               releaseDate,
			ImageSrc:           rn.ImageUrl,
			MediaLink:          rn.MediaLink,
			Text:               rn.DescriptionShort,
			LastUpdateOn:       rn.UpdatedAt.String(),
			CtaLabelOverride:   rn.CtaLabelOverride,
			CtaHrefOverride:    rn.CtaUrlOverride,
			HideCta:            rn.HideCta,
			AttentionMechanism: rn.AttentionMechanism.String(),
		})
	}

	h.log.Debug().Int("dataLength", len(res.Data)).Msg("Response data length")

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		h.log.Error().Err(err).Msg("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
