package account

import (
	"net/http"

	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// releasePageURLUpdateForm represents the form data for updating release page base URL
type releasePageURLUpdateForm struct {
	UseCustomURL       bool   `schema:"use_custom_url"`
	CustomURL          string `schema:"custom_url"`
	DisableReleasePage bool   `schema:"disable_release_page"`
}

// HandleReleasePageURLUpdate handles PATCH /settings/release-page-url
func (h *Handlers) HandleReleasePageURLUpdate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleReleasePageURLUpdate")
	ctx := r.Context()
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.deps.DB))

	orgID := ctx.Value(mw.OrgIDKey).(string)
	if orgID == "" {
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
	var updateDTO releasePageURLUpdateForm
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

	var newBaseURL *string
	if updateDTO.UseCustomURL {
		newBaseURL = &updateDTO.CustomURL
	} else {
		newBaseURL = nil
	}

	if err := widgetService.UpdateBaseURL(uuid.MustParse(orgID), newBaseURL); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating widget base URL")
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))
	if err := releasePageConfigService.UpdateDisableReleasePage(uuid.MustParse(orgID), updateDTO.DisableReleasePage); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating disable release page")
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
