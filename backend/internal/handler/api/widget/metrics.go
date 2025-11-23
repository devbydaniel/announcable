package widget

import (
	"encoding/json"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotemetrics "github.com/devbydaniel/release-notes-go/internal/domain/release-note-metrics"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createMetricRequest struct {
	ReleaseNoteID string `json:"release_note_id"`
	MetricType    string `json:"metric_type"`
	ClientID      string `json:"client_id"`
}

// HandleReleaseNoteMetricCreate creates a metric for a release note
func (h *Handlers) HandleReleaseNoteMetricCreate(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleReleaseNoteMetricCreate")
	orgService := organisation.NewService(*organisation.NewRepository(h.DB))

	// Get external org ID from URL params
	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.Log.Error().Msg("Organisation ID not found in URL")
		http.Error(w, "Organisation ID required", http.StatusBadRequest)
		return
	}

	// Get org by external ID
	org, err := orgService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting organisation")
		http.Error(w, "Error getting organisation", http.StatusInternalServerError)
		return
	}

	// Parse request body
	var req createMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error().Err(err).Msg("Error decoding request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ReleaseNoteID == "" || req.MetricType == "" || req.ClientID == "" {
		h.Log.Error().Msg("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Parse UUIDs
	releaseNoteUUID, err := uuid.Parse(req.ReleaseNoteID)
	if err != nil {
		h.Log.Error().Err(err).Msg("Invalid release note ID")
		http.Error(w, "Invalid release note ID", http.StatusBadRequest)
		return
	}

	orgUUID, err := uuid.Parse(org.ID.String())
	if err != nil {
		h.Log.Error().Err(err).Msg("Invalid organisation ID")
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	// Validate metric type
	metricType := releasenotemetrics.MetricType(req.MetricType)
	if metricType != releasenotemetrics.MetricTypeView && metricType != releasenotemetrics.MetricTypeCtaClick {
		h.Log.Error().Str("metricType", string(metricType)).Msg("Invalid metric type")
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	// Create metric
	metricsService := releasenotemetrics.NewService(releasenotemetrics.NewRepository(h.DB))
	if err := metricsService.CreateMetric(releaseNoteUUID, orgUUID, req.ClientID, metricType); err != nil {
		h.Log.Error().Err(err).Msg("Error creating metric")
		http.Error(w, "Error creating metric", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
