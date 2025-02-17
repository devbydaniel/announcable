package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	widgetconfigs "github.com/devbydaniel/release-notes-go/internal/domain/widget-configs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type serviceWidgetConfigResponseBodyWidgetConfig struct {
	OrgId                   string `json:"org_id"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	CtaText                 string `json:"cta_text"`
	WidgetType              string `json:"widget_type"`
	WidgetBorderRadius      int    `json:"widget_border_radius"`
	WidgetBorderColor       string `json:"widget_border_color"`
	WidgetBorderWidth       int    `json:"widget_border_width"`
	WidgetBgColor           string `json:"widget_bg_color"`
	WidgetFontColor         string `json:"widget_font_color"`
	ReleaseNoteBorderRadius int    `json:"release_note_border_radius"`
	ReleaseNoteBorderColor  string `json:"release_note_border_color"`
	ReleaseNoteBorderWidth  int    `json:"release_note_border_width"`
	ReleaseNoteBgColor      string `json:"release_note_bg_color"`
	ReleaseNoteFontColor    string `json:"release_note_font_color"`
	ReleasePageBaseUrl      string `json:"release_page_baseurl"`
}

type serveWidgetConfigResponseBody struct {
	Data serviceWidgetConfigResponseBodyWidgetConfig `json:"data"`
}

func (h *Handler) HandleWidgetConfigServe(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleWidgetConfigServe")
	widgetConfigService := widgetconfigs.NewService(*widgetconfigs.NewRepository(h.DB))
	organisationService := organisation.NewService(*organisation.NewRepository(h.DB))

	externalOrgId := chi.URLParam(r, "orgId")
	if externalOrgId == "" {
		h.log.Error().Msg("Org ID not found in URL")
		http.Error(w, "Error getting widget config", http.StatusBadRequest)
		return
	}

	org, err := organisationService.GetOrgByExternalId(uuid.MustParse(externalOrgId))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting org ID")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	widgetConfig, err := widgetConfigService.Get(org.ID.String())
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting widget config")
		http.Error(w, "Error getting widget config", http.StatusInternalServerError)
		return
	}

	baseUrl := config.New().BaseURL + "/s/" + externalOrgId // TODO
	if widgetConfig.ReleasePageBaseUrl != nil {
		baseUrl = *widgetConfig.ReleasePageBaseUrl
	}

	conf := serviceWidgetConfigResponseBodyWidgetConfig{
		OrgId:                   externalOrgId,
		Title:                   widgetConfig.Title,
		Description:             widgetConfig.Description,
		CtaText:                 widgetConfig.ReleaseNoteCtaText,
		WidgetType:              widgetConfig.WidgetType.String(),
		WidgetBorderRadius:      widgetConfig.WidgetBorderRadius,
		WidgetBorderColor:       widgetConfig.WidgetBorderColor,
		WidgetBorderWidth:       widgetConfig.WidgetBorderWidth,
		WidgetBgColor:           widgetConfig.WidgetBgColor,
		WidgetFontColor:         widgetConfig.WidgetTextColor,
		ReleaseNoteBorderRadius: widgetConfig.ReleaseNoteBorderRadius,
		ReleaseNoteBorderColor:  widgetConfig.ReleaseNoteBorderColor,
		ReleaseNoteBorderWidth:  widgetConfig.ReleaseNoteBorderWidth,
		ReleaseNoteBgColor:      widgetConfig.ReleaseNoteBgColor,
		ReleaseNoteFontColor:    widgetConfig.ReleaseNoteTextColor,
		ReleasePageBaseUrl:      baseUrl,
	}

	res := serveWidgetConfigResponseBody{
		Data: conf,
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
