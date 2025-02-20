package releasepageconfig

import (
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

func (s *service) Create(cfg *ReleasePageConfig, imgInput *ImageInput) (uuid.UUID, error) {
	log.Trace().Msg("Create")

	// Start a transaction
	tx := s.repo.db.StartTransaction()

	// Create image
	if imgInput != nil {
		if imgInput.ImgData != nil {
			processedImg, format, err := imgUtil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
			if err != nil {
				log.Error().Err(err).Msg("Error processing image")
				return uuid.Nil, err
			}
			path := createPath(cfg.OrganisationID.String(), format.String())
			log.Debug().Str("path", path).Msg("Creating image")
			if err := s.repo.UpdateImage(path, processedImg); err != nil {
				log.Error().Err(err).Msg("Error creating image")
				tx.Rollback()
				return uuid.Nil, err
			}
			cfg.ImagePath = path
		}
	}

	if err := s.repo.Create(cfg, tx.Tx); err != nil {
		log.Error().Err(err).Msg("Error creating widget config")
		return uuid.Nil, err
	}
	return cfg.ID, nil
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
