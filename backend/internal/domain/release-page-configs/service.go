package releasepageconfig

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/imgutil"
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

var imgProcessConfig = imgutil.ImgProcessConfig{
	MaxWidth: 1000,
	Quality:  80,
}

func createPath(orgID, format string) string {
	return orgID + "." + format
}

// NewService creates a new release page config service with the given repository.
func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) defaultConfig(orgID uuid.UUID, slug string) *ReleasePageConfig {
	log.Trace().Str("orgID", orgID.String()).Str("slug", slug).Msg("DefaultConfig")
	return &ReleasePageConfig{
		OrganisationID: orgID,
		Title:          "Release Notes",
		Description:    "See what's new",
		Slug:           slug,
		BgColor:        "#eff1f5",
		TextColor:      "#000000",
		TextColorMuted: "#4c4f69",
		BrandPosition:  BrandPositionLeft.String(),
	}
}

func (s *service) Init(orgID uuid.UUID, orgName string) (*ReleasePageConfig, error) {
	log.Trace().Str("orgID", orgID.String()).Str("orgName", orgName).Msg("Init")
	slug := s.formatSlug(orgName)
	cfg := s.defaultConfig(orgID, slug)
	if err := s.repo.Create(cfg, nil); err != nil {
		log.Error().Err(err).Msg("Error creating default config")
		return cfg, err
	}
	return cfg, nil
}

func (s *service) Get(orgID uuid.UUID) (*ReleasePageConfig, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Get")
	cfg, err := s.repo.Get(orgID)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by organisation ID")
		return nil, err
	}
	imgURL, err := s.repo.GetImageURL(cfg.ImagePath)
	if err != nil {
		log.Error().Err(err).Msg("Error getting image URL")
	}
	cfg.ImageURL = imgURL

	return cfg, nil
}

func (s *service) GetBySlug(slug string) (*ReleasePageConfig, error) {
	log.Trace().Str("slug", slug).Msg("GetBySlug")
	cfg, err := s.repo.GetBySlug(slug)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by slug")
		return nil, err
	}
	imgURL, err := s.repo.GetImageURL(cfg.ImagePath)
	if err != nil {
		log.Error().Err(err).Msg("Error getting image URL")
	}
	cfg.ImageURL = imgURL

	return cfg, nil
}

func (s *service) GetURL(orgID uuid.UUID) (string, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetURL")
	baseURL := config.New().BaseURL
	cfg, err := s.repo.Get(orgID)
	if err != nil {
		return "", err
	}
	if cfg.Slug == "" {
		// Create a new slug if it doesn't exist
		log.Debug().Msg("Slug does not exist, creating...")
		// Get organization name from the config
		if cfg.Organisation.Name == "" {
			return "", nil // Cannot create slug without org name
		}
		// Update the slug using the org name
		if err := s.UpdateSlug(orgID, cfg.Organisation.Name); err != nil {
			log.Error().Err(err).Msg("Error updating release page slug")
			return "", err
		}
		// Reload the config to get the new slug
		cfg, err = s.repo.Get(orgID)
		if err != nil {
			log.Error().Err(err).Msg("Error getting updated config")
			return "", err
		}
		if cfg.Slug == "" {
			return "", nil
		}
	}
	return baseURL + "/s/" + cfg.Slug, nil
}

func (s *service) Update(orgID uuid.UUID, cfg *ReleasePageConfig, imgInput *ImageInput) error {
	log.Trace().Msg("Update")

	// Start a transaction
	tx := s.repo.db.StartTransaction()

	// Update image
	if imgInput != nil {
		if imgInput.ShouldDeleteImage {
			if err := s.repo.DeleteImage(orgID); err != nil {
				log.Error().Err(err).Msg("Error deleting image")
				tx.Rollback()
				return err
			}
			if err := s.repo.UpdateWithNil(orgID, map[string]interface{}{"ImagePath": nil}, tx.Tx); err != nil {
				log.Error().Err(err).Msg("Error updating image path")
				tx.Rollback()
				return err
			}
		} else if imgInput.ImgData != nil {
			processedImg, format, err := imgutil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
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

	if err := s.repo.Update(orgID, cfg, tx.Tx); err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) UpdateSlug(orgID uuid.UUID, orgName string) error {
	log.Trace().Str("orgID", orgID.String()).Str("orgName", orgName).Msg("UpdateSlug")
	slug := s.formatSlug(orgName)
	return s.repo.UpdateWithNil(orgID, map[string]interface{}{"Slug": slug}, nil)
}

func (s *service) EditSlugAsAdmin(orgID uuid.UUID, slug string) error {
	log.Trace().Str("orgID", orgID.String()).Str("slug", slug).Msg("EditSlugAsAdmin")
	// Validate slug format
	if len(slug) <= 0 {
		return errors.New("Slug is required")
	}
	if len(slug) > 100 {
		return errors.New("Slug is too long")
	}
	// only allow alphanumeric characters and hyphens
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(slug) {
		return errors.New("Slug can only contain alphanumeric characters and hyphens")
	}
	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Check if slug is already in use by another organization
	existingConfig, err := s.repo.GetBySlug(slug)
	if err == nil && existingConfig.OrganisationID != orgID {
		return errors.New("Slug is already in use by another organization")
	}

	return s.repo.UpdateWithNil(orgID, map[string]interface{}{"Slug": slug}, nil)
}

func (s *service) UpdateDisableReleasePage(orgID uuid.UUID, disabled bool) error {
	log.Trace().Str("orgID", orgID.String()).Bool("disabled", disabled).Msg("UpdateDisableReleasePage")
	return s.repo.UpdateWithNil(orgID, map[string]interface{}{"DisableReleasePage": disabled}, nil)
}

func (s *service) formatSlug(orgName string) string {
	log.Trace().Str("orgName", orgName).Msg("GetSlug")
	// make lowercase and url friendly
	slug := strings.ToLower(orgName)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = url.PathEscape(slug)
	return slug
}
