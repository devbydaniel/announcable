package handler

import (
	"net/http"

	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type widgetUpdateForm struct {
	Title                   string `schema:"title" validate:"required"`
	Description             string `schema:"description" validate:"required"`
	WidgetType              string `schema:"widget_type" validate:"required"`
	WidgetBorderRadius      int    `schema:"widget_border_radius" validate:"gte=0"`
	WidgetBorderColor       string `schema:"widget_border_color" validate:"required"`
	WidgetBorderWidth       int    `schema:"widget_border_width" validate:"gte=0"`
	WidgetBgColor           string `schema:"widget_background_color" validate:"required"`
	WidgetTextColor         string `schema:"widget_text_color" validate:"required"`
	ReleaseNoteBorderRadius int    `schema:"release_note_border_radius" validate:"gte=0"`
	ReleaseNoteBorderColor  string `schema:"release_note_border_color" validate:"required"`
	ReleaseNoteBorderWidth  int    `schema:"release_note_border_width" validate:"gte=0"`
	ReleaseNoteBgColor      string `schema:"release_note_background_color" validate:"required"`
	ReleaseNoteTextColor    string `schema:"release_note_text_color" validate:"required"`
	ReleaseNoteCtaText      string `schema:"release_note_cta_text" validate:"required"`
	EnableLikes             string `schema:"enable_likes"`
	LikeButtonText          string `schema:"like_button_text"`
	UnlikeButtonText        string `schema:"unlike_button_text"`
}

func (h *Handler) HandleWidgetUpdate(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleWidgetUpdate")
	ctx := r.Context()
	widgetService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))

	userId := ctx.Value(mw.UserIDKey).(string)
	if userId == "" {
		h.log.Error().Msg("User ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
	}

	// parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO widgetUpdateForm
	if err := h.decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating widget config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Additional validation for like button text fields
	isLikesEnabled := updateDTO.EnableLikes == "on"
	if isLikesEnabled {
		if updateDTO.LikeButtonText == "" {
			http.Error(w, "Like button text is required when likes are enabled", http.StatusBadRequest)
			return
		}
		if updateDTO.UnlikeButtonText == "" {
			http.Error(w, "Unlike button text is required when likes are enabled", http.StatusBadRequest)
			return
		}
	} else {
		// Clear button text fields when likes are disabled
		updateDTO.LikeButtonText = ""
		updateDTO.UnlikeButtonText = ""
	}

	widgetConfig := &widgetconfigs.WidgetConfig{
		OrganisationID:          uuid.MustParse(orgId),
		Title:                   updateDTO.Title,
		Description:             updateDTO.Description,
		WidgetType:              widgetconfigs.WidgetType(updateDTO.WidgetType),
		WidgetBorderRadius:      updateDTO.WidgetBorderRadius,
		WidgetBorderColor:       updateDTO.WidgetBorderColor,
		WidgetBorderWidth:       updateDTO.WidgetBorderWidth,
		WidgetBgColor:           updateDTO.WidgetBgColor,
		WidgetTextColor:         updateDTO.WidgetTextColor,
		ReleaseNoteBorderRadius: updateDTO.ReleaseNoteBorderRadius,
		ReleaseNoteBorderColor:  updateDTO.ReleaseNoteBorderColor,
		ReleaseNoteBorderWidth:  updateDTO.ReleaseNoteBorderWidth,
		ReleaseNoteBgColor:      updateDTO.ReleaseNoteBgColor,
		ReleaseNoteTextColor:    updateDTO.ReleaseNoteTextColor,
		ReleaseNoteCtaText:      updateDTO.ReleaseNoteCtaText,
		EnableLikes:             isLikesEnabled,
		LikeButtonText:          updateDTO.LikeButtonText,
		UnlikeButtonText:        updateDTO.UnlikeButtonText,
	}
	h.log.Debug().Interface("widget config", widgetConfig).Msg("Widget config to update")

	if err := widgetService.Update(uuid.MustParse(orgId), widgetConfig); err != nil {
		h.log.Error().Err(err).Msg("Error updating widget config")
		http.Error(w, "Error updating widget config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
