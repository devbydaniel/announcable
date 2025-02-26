package releasepageconfig

import (
	"net/url"
	"strings"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/imgUtil"
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

var imgProcessConfig = imgUtil.ImgProcessConfig{
	MaxWidth: 1000,
	Quality:  80,
}

func createPath(orgId, format string) string {
	return orgId + "." + format
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) defaultConfig(orgId uuid.UUID, slug string) *ReleasePageConfig {
	log.Trace().Str("orgId", orgId.String()).Str("slug", slug).Msg("DefaultConfig")
	return &ReleasePageConfig{
		OrganisationID: orgId,
		Title:          "Release Notes",
		Description:    "See what's new",
		Slug:           slug,
		BgColor:        "#eff1f5",
		TextColor:      "#000000",
		TextColorMuted: "#4c4f69",
		BrandPosition:  BrandPositionLeft.String(),
	}
}

func (s *service) Init(orgId uuid.UUID, orgName string) (*ReleasePageConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Str("orgName", orgName).Msg("Init")
	slug := s.formatSlug(orgName)
	cfg := s.defaultConfig(orgId, slug)
	if err := s.repo.Create(cfg, nil); err != nil {
		log.Error().Err(err).Msg("Error creating default config")
		return cfg, err
	}
	return cfg, nil
}

func (s *service) Get(orgId uuid.UUID) (*ReleasePageConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("Get")
	cfg, err := s.repo.Get(orgId)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by organisation ID")
		return nil, err
	}
	imgUrl, err := s.repo.GetImageUrl(cfg.ImagePath)
	if err != nil {
		log.Error().Err(err).Msg("Error getting image URL")
	}
	cfg.ImageUrl = imgUrl

	return cfg, nil
}

func (s *service) GetBySlug(slug string) (*ReleasePageConfig, error) {
	log.Trace().Str("slug", slug).Msg("GetBySlug")
	cfg, err := s.repo.GetBySlug(slug)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by slug")
		return nil, err
	}
	imgUrl, err := s.repo.GetImageUrl(cfg.ImagePath)
	if err != nil {
		log.Error().Err(err).Msg("Error getting image URL")
	}
	cfg.ImageUrl = imgUrl

	return cfg, nil
}

func (s *service) GetUrl(orgId uuid.UUID) (string, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetUrl")
	baseUrl := config.New().BaseURL
	cfg, err := s.repo.Get(orgId)
	if err != nil {
		return "", err
	}
	if cfg.Slug == "" {
		return "", nil
	}
	return baseUrl + "/s/" + cfg.Slug, nil
}

func (s *service) Update(orgId uuid.UUID, cfg *ReleasePageConfig, imgInput *ImageInput) error {
	log.Trace().Msg("Update")

	// Start a transaction
	tx := s.repo.db.StartTransaction()

	// Update image
	if imgInput != nil {
		if imgInput.ShoudDeleteImage {
			if err := s.repo.DeleteImage(orgId); err != nil {
				log.Error().Err(err).Msg("Error deleting image")
				tx.Rollback()
				return err
			}
			if err := s.repo.UpdateWithNil(orgId, map[string]interface{}{"ImagePath": nil}, tx.Tx); err != nil {
				log.Error().Err(err).Msg("Error updating image path")
				tx.Rollback()
				return err
			}
		} else if imgInput.ImgData != nil {
			processedImg, format, err := imgUtil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
			if err != nil {
				log.Error().Err(err).Msg("Error processing image")
				tx.Rollback()
				return err
			}
			path := createPath(cfg.OrganisationID.String(), format.String())
			log.Debug().Str("path", path).Msg("Updating image")
			if err := s.repo.UpdateImage(path, processedImg); err != nil {
				log.Error().Err(err).Msg("Error updating image")
				tx.Rollback()
				return err
			}
			cfg.ImagePath = path
		}
	}

	if err := s.repo.Update(orgId, cfg, tx.Tx); err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) UpdateSlug(orgId uuid.UUID, orgName string) error {
	log.Trace().Str("orgId", orgId.String()).Str("orgName", orgName).Msg("UpdateSlug")
	slug := s.formatSlug(orgName)
	return s.repo.UpdateWithNil(orgId, map[string]interface{}{"Slug": slug}, nil)
}

func (s *service) formatSlug(orgName string) string {
	log.Trace().Str("orgName", orgName).Msg("GetSlug")
	// make lowercase and url friendly
	slug := strings.ToLower(orgName)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = url.PathEscape(slug)
	return slug
}
