package account

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/google/uuid"
)

// HandleWidgetIdRegenerate handles PATCH /settings/widget-id
func (h *Handlers) HandleWidgetIdRegenerate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleWidgetIdRegenerate")
	ctx := r.Context()
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	organisationService := organisation.NewService(*organisation.NewRepository(h.deps.DB))

	_, err := organisationService.RegenerateExternalId(uuid.MustParse(orgId))
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error regenerating widget external ID")
		http.Error(w, "Error regenerating widget external ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
