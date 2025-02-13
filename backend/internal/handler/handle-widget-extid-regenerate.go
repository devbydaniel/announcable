package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/google/uuid"
)

func (h *Handler) HandleOrgExternalIdRegenerate(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleWidgetExternalIdRegenerate")
	ctx := r.Context()
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))

	_, err := organisationService.RegenerateExternalId(uuid.MustParse(orgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error regenerating widget external ID")
		http.Error(w, "Error regenerating widget external ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
	return
}
