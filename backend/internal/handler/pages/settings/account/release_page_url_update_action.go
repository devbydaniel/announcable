package account

import (
	"net/http"

	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// releasePageUrlUpdateForm represents the form data for updating release page base URL
type releasePageUrlUpdateForm struct {
	UseCustomUrl bool   `schema:"use_custom_url"`
	CustomUrl    string `schema:"custom_url"`
}

// HandleReleasePageUrlUpdate handles PATCH /settings/release-page-url
func (h *Handlers) HandleReleasePageUrlUpdate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleReleasePageUrlUpdate")
	ctx := r.Context()
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating release page URL", http.StatusInternalServerError)
		return
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating release page URL", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO releasePageUrlUpdateForm
	if err := h.deps.Decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating release page URL", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating release page URL", http.StatusBadRequest)
		return
	}

	var newBaseUrl *string
	if updateDTO.UseCustomUrl {
		newBaseUrl = &updateDTO.CustomUrl
	} else {
		newBaseUrl = nil
	}

	widgetService.UpdateBaseUrl(uuid.MustParse(orgId), newBaseUrl)

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
