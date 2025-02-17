package handler

import (
	"net/http"

	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-playground/validator"
)

type lpBaseUrlUpdateForm struct {
	UseCustomUrl bool   `schema:"use_custom_url"`
	CustomUrl    string `schema:"custom_url"`
}

func (h *Handler) HandleReleasePageBaseUrlUpdate(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleReleasePageBaseUrlUpdate")
	ctx := r.Context()
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO lpBaseUrlUpdateForm
	if err := h.decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	var newBaseUrl *string
	if updateDTO.UseCustomUrl {
		newBaseUrl = &updateDTO.CustomUrl
	} else {
		newBaseUrl = nil
	}

	widgetService.UpdateBaseUrl(orgId, newBaseUrl)

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
	return
}
