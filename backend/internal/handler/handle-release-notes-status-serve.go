package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
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

func (h *Handler) HandleReleaseNotesStatusServe(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleaseNotesServe")
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))
	forWidgetOrWebsite := r.URL.Query().Get("for")
	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.log.Error().Msg("Org ID not found in URL")
		http.Error(w, "Error getting release notes status", http.StatusBadRequest)
		return
	}

	org, err := organisationService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting org ID")
		http.Error(w, "Error getting release notes status", http.StatusInternalServerError)
		return
	}
	filters := map[string]interface{}{}
	if forWidgetOrWebsite == "widget" {
		filters["hide_on_widget"] = false
	} else if forWidgetOrWebsite == "website" {
		filters["hide_on_release_page"] = false
	}
	releaseNotesStatus, err := releaseNotesService.GetStatus(org.ID.String(), filters)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting release notes status")
		http.Error(w, "Error getting release notes status", http.StatusInternalServerError)
		return
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
		h.log.Error().Err(err).Msg("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
