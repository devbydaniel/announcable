package releasenotes

import (
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

func createImgPath(orgID, randomID, format string) string {
	return orgID + "/" + randomID + "." + format
}

// NewService creates a new release notes service with the given repository.
func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) Create(rn *ReleaseNote, imgInput *ImageInput) (uuid.UUID, error) {
	log.Trace().Msg("Create")
	// Start a transaction
	tx := s.repo.db.StartTransaction()

	// Create release note
	id, err := s.repo.Create(rn, tx.Tx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating release note")
		tx.Rollback()
		return uuid.Nil, err
	}

	// Create image
	if imgInput != nil {
		if imgInput.ImgData != nil {
			processedImg, format, err := imgutil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
			if err != nil {
				log.Error().Err(err).Msg("Error processing image")
				tx.Rollback()
				return uuid.Nil, err
			}
			randomID, err := uuid.NewRandom()
			if err != nil {
				log.Error().Err(err).Msg("Error generating random UUID")
				tx.Rollback()
				return uuid.Nil, err
			}
			imgPath := createImgPath(rn.OrganisationID.String(), randomID.String(), format.String())
			log.Debug().Str("path", imgPath).Msg("Creating image")
			if err := s.repo.UpdateImage(id, processedImg, imgPath, tx.Tx); err != nil {
				log.Error().Err(err).Msg("Error creating image")
				tx.Rollback()
				return uuid.Nil, err
			}
			rn.ImagePath = imgPath
		}
	}

	// Commit the transaction
	tx.Commit()
	return id, nil
}

func (s *service) GetAll(orgID string, page, pageSize int) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgID", orgID).Msg("GetAll")
	rns, err := s.repo.FindAll(orgID, page, pageSize, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error finding release notes by organisation ID")
		return nil, err
	}
	// adjust release date format
	for _, rn := range rns.Items {
		if rn.ReleaseDate != nil {
			rd := (*rn.ReleaseDate)[:10]
			rn.ReleaseDate = &rd
		}
	}
	return rns, nil
}

func (s *service) GetAllWithImgURL(orgID string, page, pageSize int, filters map[string]interface{}) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgID", orgID).Msg("GetAllWithImgURL")
	rns, err := s.repo.FindAll(orgID, page, pageSize, filters)
	if err != nil {
		log.Error().Err(err).Msg("Error finding release notes by organisation ID")
		return nil, err
	}
	// adjust release date format and get image
	for _, rn := range rns.Items {
		if rn.ReleaseDate != nil {
			rd := (*rn.ReleaseDate)[:10]
			rn.ReleaseDate = &rd
		}
		if rn.ImagePath != "" {
			imgURL, err := s.repo.GetImageURL(rn.ImagePath)
			if err != nil {
				log.Error().Err(err).Msg("Error getting image URL")
			} else {
				rn.ImageURL = imgURL
			}
		}
	}
	return rns, nil
}

func (s *service) GetStatus(orgID string, filters map[string]interface{}) ([]*ReleaseNoteStatus, error) {
	log.Trace().Str("orgID", orgID).Msg("GetStatus")
	return s.repo.GetStatus(orgID, filters)
}

func (s *service) GetOne(id, orgID string) (*ReleaseNote, error) {
	log.Trace().Msg("GetByID")

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing UUID")
		return nil, err
	}

	rn, err := s.repo.FindOne(uuid)
	if err != nil {
		log.Error().Err(err).Msg("Error finding release note by ID")
		return nil, err
	}

	// ReleaseDate includes time, so we need to format it
	if rn.ReleaseDate != nil {
		rd := (*rn.ReleaseDate)[:10]
		rn.ReleaseDate = &rd
	}

	// Get image URL from path if it exists
	if rn.ImagePath != "" {
		imgURL, err := s.repo.GetImageURL(rn.ImagePath)
		if err != nil {
			log.Error().Err(err).Msg("Error getting image URL")
		} else {
			rn.ImageURL = imgURL
		}
	}
	return rn, nil
}

func (s *service) Update(id uuid.UUID, rn *ReleaseNote, imgInput *ImageInput) error {
	log.Trace().Msg("UpdateWithImg")
	// Start a transaction
	tx := s.repo.db.StartTransaction()

	var path string
	// Update image
	if imgInput != nil {
		if imgInput.ShouldDeleteImage {
			if err := s.repo.DeleteImage(id, tx.Tx); err != nil {
				log.Error().Err(err).Msg("Error deleting image")
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
			randID, err := uuid.NewRandom()
			if err != nil {
				log.Error().Err(err).Msg("Error generating random UUID")
				tx.Rollback()
				return err
			}
			path = createImgPath(rn.OrganisationID.String(), randID.String(), format.String())
			log.Debug().Str("path", path).Msg("Updating image")
			if err := s.repo.UpdateImage(id, processedImg, path, tx.Tx); err != nil {
				log.Error().Err(err).Msg("Error updating image")
				tx.Rollback()
				return err
			}
		}
	}

	// Update release note data
	log.Debug().Interface("rn", rn).Msg("Updating release note")
	if err := s.repo.Update(id, rn, tx.Tx); err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		if imgInput != nil && imgInput.ImgData != nil {
			_ = s.repo.DeleteImage(id, tx.Tx)
		}
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) ChangePublishedStatus(id uuid.UUID, published bool) error {
	log.Trace().Bool("published", published).Msg("ChangePublishedStatus")
	if err := s.repo.UpdateWithNil(id, map[string]interface{}{"IsPublished": published}, nil); err != nil {
		log.Error().Err(err).Msg("Error updating published status")
		return err
	}
	return nil
}

func (s *service) Delete(id uuid.UUID) error {
	log.Trace().Msg("Delete")
	return s.repo.Delete(id, nil)
}

func (s *service) GetCount(orgID uuid.UUID) (int64, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetCount")
	return s.repo.GetCount(orgID)
}
