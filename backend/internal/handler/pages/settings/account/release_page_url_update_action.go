package account

import (
	"net/http"

	releasepageconfig "github.com/devbydaniel/announcable/internal/domain/release-page-configs"
	widgetconfigs "github.com/devbydaniel/announcable/internal/domain/widget-configs"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// releasePageUrlUpdateForm represents the form data for updating release page base URL
type releasePageUrlUpdateForm struct {
	UseCustomUrl       bool   `schema:"use_custom_url"`
	CustomUrl          string `schema:"custom_url"`
	DisableReleasePage bool   `schema:"disable_release_page"`
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

	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.deps.DB, h.deps.ObjStore))
	if err := releasePageConfigService.UpdateDisableReleasePage(uuid.MustParse(orgId), updateDTO.DisableReleasePage); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating disable release page")
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
