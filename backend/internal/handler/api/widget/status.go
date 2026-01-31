package widget

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type releaseNotesStatusResponseBodyData struct {
	LastUpdateOn       string `json:"last_update_on"`
	AttentionMechanism string `json:"attention_mechanism"`
}

type releaseNotesStatusResponseBody struct {
	Data []releaseNotesStatusResponseBodyData `json:"data"`
}

// MemCacheReleaseNotesStatus is the in-memory cache for release notes status responses.
var MemCacheReleaseNotesStatus = memcache.New(5*time.Minute, 10*time.Minute)

// HandleReleaseNotesStatusServe serves release notes status for the widget
func (h *Handlers) HandleReleaseNotesStatusServe(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleReleaseNotesStatusServe")
	externalOrgID := chi.URLParam(r, "orgID")
	if externalOrgID == "" {
		h.Log.Error().Msg("Org ID not found in URL")
		http.Error(w, "Error getting release notes status", http.StatusBadRequest)
		return
	}

	var releaseNotesStatus []*releasenotes.ReleaseNoteStatus
	status, found := MemCacheReleaseNotesStatus.Get(externalOrgID)
	if found {
		h.Log.Trace().Msg("Found in cache")
		releaseNotesStatus = status.([]*releasenotes.ReleaseNoteStatus)
	} else {
		organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
		releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
		forWidgetOrWebsite := r.URL.Query().Get("for")
		h.Log.Trace().Msg("Not found in cache")
		org, err := organisationService.GetOrgByExternalID(uuid.MustParse(externalOrgID))
		if err != nil {
			h.Log.Error().Err(err).Msg("Error getting org ID")
			http.Error(w, "Error getting release notes status", http.StatusInternalServerError)
			return
		}
		filters := map[string]interface{}{}
		if forWidgetOrWebsite == "widget" {
			filters["hide_on_widget"] = false
		} else if forWidgetOrWebsite == "website" {
			filters["hide_on_release_page"] = false
		}
		status, err := releaseNotesService.GetStatus(org.ID.String(), filters)
		if err != nil {
			h.Log.Error().Err(err).Msg("Error getting release notes status")
			http.Error(w, "Error getting release notes status", http.StatusInternalServerError)
			return
		}
		MemCacheReleaseNotesStatus.Set(externalOrgID, status, 5*time.Minute)
		releaseNotesStatus = status
	}

	var res releaseNotesStatusResponseBody
	for _, rn := range releaseNotesStatus {
		res.Data = append(res.Data, releaseNotesStatusResponseBodyData{
			LastUpdateOn:       rn.UpdatedAt,
			AttentionMechanism: rn.AttentionMechanism,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		h.Log.Error().Err(err).Msg("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
