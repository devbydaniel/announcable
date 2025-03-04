package handler

import (
	"net/http"

	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/devbydaniel/release-notes-go/internal/imgUtil"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type releaseNoteUpdateForm struct {
	// The request also contains a field for the image,
	// but gorilla cannot parse it as it is a file in a multipart form.
	Title               string `schema:"title" validate:"required"`
	ShouldDeleteImage   bool   `schema:"delete_image"`
	DescriptionShort    string `schema:"description_short" validate:"required"`
	TextWebsiteOverride string `schema:"text_website_override"`
	DescriptionLong     string `schema:"description_long"`
	ReleaseDate         string `schema:"release_date"`
	HideCta             bool   `schema:"hide_cta"`
	OverrideCtaLabel    bool   `schema:"override_cta_label"`
	CtaLabelOverride    string `schema:"cta_label_override"`
	OverrideCtaUrl      bool   `schema:"override_cta_url"`
	CtaUrlOverride      string `schema:"cta_url_override"`
	AttentionMechanism  string `schema:"attention_mechanism"`
	HideOnWidget        string `schema:"hide_on_widget"`
	HideOnReleasePage   string `schema:"hide_on_release_page"`
}

func (h *Handler) HandleReleaseNoteUpdate(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleNewReleaseNote")
	ctx := r.Context()
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.DB, h.ObjStore))

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Error().Err(err).Msg("Error parsing ID")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error updating release note", http.StatusBadRequest)
		return
	}

	// decode form
	var updateDTO releaseNoteUpdateForm
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

	// get image
	img, imgHeader, err := r.FormFile("image")
	if err != nil {
		h.log.Debug().Err(err).Msg("No image uploaded")
	}
	h.log.Debug().Interface("img", img).Msg("img")

	// prepare models
	var imgInput *releasenotes.ImageInput
	if img != nil {
		ok := imgUtil.VerifyImageType(img)
		if !ok {
			h.log.Error().Msg("Invalid image type")
			http.Error(w, "Error updating release note", http.StatusBadRequest)
			return
		}
		imgInput = &releasenotes.ImageInput{
			ShouldDeleteImage: false,
			ImgData:           img,
			Format:            imgHeader.Header.Get("Content-Type"),
		}
	} else {
		imgInput = &releasenotes.ImageInput{
			ShouldDeleteImage: updateDTO.ShouldDeleteImage,
			ImgData:           nil,
			Format:            "",
		}
	}
	h.log.Debug().Interface("imgInput", imgInput).Msg("ImageInput")
	h.log.Debug().Interface("updateDTO", updateDTO).Msg("updateDTO")

	releaseNote := &releasenotes.ReleaseNote{
		OrganisationID:     uuid.MustParse(orgId),
		Title:              updateDTO.Title,
		DescriptionShort:   updateDTO.DescriptionShort,
		HideCta:            updateDTO.HideCta,
		AttentionMechanism: releasenotes.AttentionMechanism(updateDTO.AttentionMechanism),
		LastUpdatedBy:      uuid.MustParse(userId),
		HideOnWidget:       updateDTO.HideOnWidget == "on",
		HideOnReleasePage:  updateDTO.HideOnReleasePage == "on",
	}
	if updateDTO.TextWebsiteOverride == "on" {
		releaseNote.DescriptionLong = updateDTO.DescriptionLong
	} else {
		releaseNote.DescriptionLong = ""
	}
	if !updateDTO.HideCta && updateDTO.OverrideCtaLabel {
		releaseNote.CtaLabelOverride = updateDTO.CtaLabelOverride
	} else {
		releaseNote.CtaLabelOverride = ""
	}
	if !updateDTO.HideCta && updateDTO.OverrideCtaUrl {
		releaseNote.CtaUrlOverride = updateDTO.CtaUrlOverride
	} else {
		releaseNote.CtaUrlOverride = ""
	}
	if updateDTO.ReleaseDate != "" {
		releaseNote.ReleaseDate = &updateDTO.ReleaseDate
	} else {
		releaseNote.ReleaseDate = nil
	}
	h.log.Debug().Interface("releaseNote", releaseNote).Msg("ReleaseNote to update")

	if err := releaseNotesService.Update(id, releaseNote, imgInput); err != nil {
		h.log.Error().Err(err).Msg("Error updating release note")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "custom:submit-success")
	w.WriteHeader(http.StatusOK)
}
