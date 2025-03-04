package handler

import (
	"net/http"

	releasepageconfig "github.com/devbydaniel/release-notes-go/internal/domain/release-page-configs"
	"github.com/devbydaniel/release-notes-go/internal/imgUtil"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type lpConfigUpdateForm struct {
	// The request also contains a field for the image,
	// but gorilla cannot parse it as it is a file in a multipart form.
	ShouldDeleteImage bool   `schema:"delete_image"`
	BrandPosition     string `schema:"brand_position" validate:"required"`
	Title             string `schema:"title" validate:"required"`
	Description       string `schema:"description" validate:"required"`
	BgColor           string `schema:"bg_color" validate:"required"`
	TextColor         string `schema:"text_color" validate:"required"`
	TextColorMuted    string `schema:"text_color_muted" validate:"required"`
	BackLinkLabel     string `schema:"back_link_label"`
	BackLinkUrl       string `schema:"back_link_url"`
}

func (h *Handler) HandleReleasePageConfigUpdate(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleLpConfigUpdate")
	ctx := r.Context()
	releasePageConfigService := releasepageconfig.NewService(*releasepageconfig.NewRepository(h.DB, h.ObjStore))

	orgId := ctx.Value(mw.OrgIDKey).(string)
	if orgId == "" {
		h.log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error updating landing page config", http.StatusInternalServerError)
	}

	// parse form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating landing page config", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO lpConfigUpdateForm
	if err := h.decoder.Decode(&updateDTO, r.PostForm); err != nil {
		h.log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error updating landing page config", http.StatusBadRequest)
		return
	}

	// validate form
	validate := validator.New()
	if err := validate.Struct(updateDTO); err != nil {
		h.log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Error updating landing page config", http.StatusBadRequest)
		return
	}

	// get image
	img, imgHeader, err := r.FormFile("logo")
	if err != nil {
		h.log.Debug().Err(err).Msg("No image uploaded")
	}

	// prepare models
	var imgInput *releasepageconfig.ImageInput
	if img != nil {
		ok := imgUtil.VerifyImageType(img)
		if !ok {
			h.log.Error().Msg("Invalid image type")
			http.Error(w, "Error updating release note", http.StatusBadRequest)
			return
		}
		imgInput = &releasepageconfig.ImageInput{
			ShouldDeleteImage: false,
			ImgData:           img,
			Format:            imgHeader.Header.Get("Content-Type"),
		}
	} else {
		imgInput = &releasepageconfig.ImageInput{
			ShouldDeleteImage: updateDTO.ShouldDeleteImage,
			ImgData:           nil,
			Format:            "",
		}
	}
	h.log.Debug().Interface("imgInput", imgInput).Msg("ImageInput")

	lpConfig := &releasepageconfig.ReleasePageConfig{
		OrganisationID: uuid.MustParse(orgId),
		BrandPosition:  updateDTO.BrandPosition,
		Title:          updateDTO.Title,
		Description:    updateDTO.Description,
		BgColor:        updateDTO.BgColor,
		TextColor:      updateDTO.TextColor,
		TextColorMuted: updateDTO.TextColorMuted,
		BackLinkLabel:  updateDTO.BackLinkLabel,
		BackLinkUrl:    updateDTO.BackLinkUrl,
	}
	h.log.Debug().Interface("lpconfig", lpConfig).Msg("Landing page config to update")

	if err := releasePageConfigService.Update(uuid.MustParse(orgId), lpConfig, imgInput); err != nil {
		h.log.Error().Err(err).Msg("Error updating landing page config")
		http.Error(w, "Error updating landing page config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
