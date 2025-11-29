package create

import (
	"fmt"
	"net/http"
	"net/url"

	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/devbydaniel/announcable/internal/imgUtil"
	mw "github.com/devbydaniel/announcable/internal/middleware"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

type releaseNoteCreateForm struct {
	// The request also contains a field for the image,
	// but gorilla cannot parse it as it is a file in a multipart form.
	Title               string `schema:"title" validate:"required"`
	DescriptionShort    string `schema:"description_short" validate:"required"`
	DeleteImage         bool   `schema:"delete_image"`
	TextWebsiteOverride string `schema:"text_website_override"`
	DescriptionLong     string `schema:"description_long"`
	ReleaseDate         string `schema:"release_date"`
	OverrideCtaLabel    bool   `schema:"override_cta_label"`
	CtaLabelOverride    string `schema:"cta_label_override"`
	OverrideCtaUrl      bool   `schema:"override_cta_url"`
	CtaUrlOverride      string `schema:"cta_url_override"`
	HideCta             bool   `schema:"hide_cta"`
	AttentionMechanism  string `schema:"attention_mechanism"`
	HideOnWidget        bool   `schema:"hide_on_widget"`
	HideOnReleasePage   bool   `schema:"hide_on_release_page"`
	MediaType           string `schema:"media_type"`
	MediaLink           string `schema:"media_link"`
}

// HandleReleaseNoteCreate handles POST /release-notes/
func (h *Handlers) HandleReleaseNoteCreate(w http.ResponseWriter, r *http.Request) {
	h.deps.Log.Trace().Msg("HandleReleaseNoteCreate")

	ctx := r.Context()
	releaseNotesService := releasenotes.NewService(*releasenotes.NewRepository(h.deps.DB, h.deps.ObjStore))

	// extract organisation ID from context
	orgId, ok := ctx.Value(mw.OrgIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("Organisation ID not found in context")
		http.Error(w, "Error while authenticating", http.StatusInternalServerError)
		return
	}
	userId, ok := ctx.Value(mw.UserIDKey).(string)
	if !ok {
		h.deps.Log.Error().Msg("User ID not found in context")
		http.Error(w, "Error while authenticating", http.StatusInternalServerError)
		return
	}

	// parse form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// decode form
	var createDTO releaseNoteCreateForm
	if err := schema.NewDecoder().Decode(&createDTO, r.PostForm); err != nil {
		h.deps.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}
	h.deps.Log.Debug().Interface("createDTO", createDTO).Msg("createDTO")

	// validate form
	validate := validator.New()
	if err := validate.Struct(createDTO); err != nil {
		h.deps.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	// Prepare the image input based on media type
	var imgInput *releasenotes.ImageInput
	if createDTO.MediaType == "image" {
		// Only process image if image type is selected
		img, imgHeader, err := r.FormFile("image")
		if err != nil {
			h.deps.Log.Debug().Err(err).Msg("No image uploaded")
		}

		if img != nil {
			ok := imgUtil.VerifyImageType(img)
			if !ok {
				h.deps.Log.Error().Msg("Invalid image type")
				http.Error(w, "Error updating release note", http.StatusBadRequest)
				return
			}
			imgInput = &releasenotes.ImageInput{
				ShouldDeleteImage: false,
				ImgData:           img,
				Format:            imgHeader.Header.Get("Content-Type"),
			}
		} else if createDTO.DeleteImage {
			imgInput = &releasenotes.ImageInput{
				ShouldDeleteImage: true,
			}
		}
	}

	releaseNote := &releasenotes.ReleaseNote{
		OrganisationID:     uuid.MustParse(orgId),
		CreatedBy:          uuid.MustParse(userId),
		Title:              createDTO.Title,
		DescriptionShort:   createDTO.DescriptionShort,
		HideCta:            createDTO.HideCta,
		AttentionMechanism: releasenotes.AttentionMechanism(createDTO.AttentionMechanism),
		LastUpdatedBy:      uuid.MustParse(userId),
		HideOnWidget:       createDTO.HideOnWidget,
		HideOnReleasePage:  createDTO.HideOnReleasePage,
	}

	// Handle media type switching
	switch createDTO.MediaType {
	case "embed":
		// If switching to embed, clear image path, set media link, and always delete any existing image
		releaseNote.ImagePath = ""
		releaseNote.MediaLink = createDTO.MediaLink
		imgInput = &releasenotes.ImageInput{
			ShouldDeleteImage: true,
		}
	case "image":
		// If switching to image, clear media link
		releaseNote.MediaLink = ""
		// Image handling is done above with imgInput
	}

	if createDTO.TextWebsiteOverride == "on" {
		releaseNote.DescriptionLong = createDTO.DescriptionLong
	} else {
		releaseNote.DescriptionLong = ""
	}
	if !createDTO.HideCta && createDTO.OverrideCtaLabel {
		releaseNote.CtaLabelOverride = createDTO.CtaLabelOverride
	} else {
		releaseNote.CtaLabelOverride = ""
	}
	if !createDTO.HideCta && createDTO.OverrideCtaUrl {
		releaseNote.CtaUrlOverride = createDTO.CtaUrlOverride
	} else {
		releaseNote.CtaUrlOverride = ""
	}
	if createDTO.ReleaseDate != "" {
		releaseNote.ReleaseDate = &createDTO.ReleaseDate
	} else {
		releaseNote.ReleaseDate = nil
	}

	id, err := releaseNotesService.Create(releaseNote, imgInput)
	if err != nil {
		h.deps.Log.Error().Err(err).Msg("Error updating release note")
		http.Error(w, "Error updating release note", http.StatusInternalServerError)
		return
	}

	successMsg := "release note created"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := fmt.Sprintf("/release-notes/%s?success=%s", id, escapedMsg)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
}
