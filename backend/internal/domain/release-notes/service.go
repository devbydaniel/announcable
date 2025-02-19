package releasenotes

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
	Format:   imgUtil.JPEG,
}

func createImgPath(orgId, rnId string) string {
	return orgId + "/" + rnId + "." + imgProcessConfig.Format.ToEncodedFormat().String()
}

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
			processedImg, err := imgUtil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
			if err != nil {
				log.Error().Err(err).Msg("Error processing image")
				tx.Rollback()
				return uuid.Nil, err
			}
			path := createImgPath(rn.OrganisationID.String(), id.String())
			log.Debug().Str("path", path).Msg("Creating image")
			if err := s.repo.UpdateImage(path, processedImg); err != nil {
				log.Error().Err(err).Msg("Error creating image")
				tx.Rollback()
				return uuid.Nil, err
			}
		}
	}

	// Commit the transaction
	tx.Commit()
	return id, nil
}

func (s *service) GetAll(orgId string, page, pageSize int) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgId", orgId).Msg("GetAll")
	rns, err := s.repo.FindAll(orgId, page, pageSize, nil)
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

func (s *service) GetAllWithImgUrl(orgId string, page, pageSize int) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgId", orgId).Msg("GetAllWithImgUrl")
	rns, err := s.repo.FindAll(orgId, page, pageSize, nil)
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
		imgUrl, err := s.repo.GetImageUrl(createImgPath(orgId, rn.ID.String()))
		if err != nil {
			log.Error().Err(err).Msg("Error getting image URL")
		}
		rn.ImageUrl = imgUrl
	}
	return rns, nil
}

func (s *service) GetStatus(orgId string) ([]*ReleaseNoteStatus, error) {
	log.Trace().Str("orgId", orgId).Msg("GetStatus")
	return s.repo.GetStatus(orgId)
}

func (s *service) GetAllPublished(orgId string, page, pageSize int) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgId", orgId).Msg("GetAllPublished")
	filter := map[string]interface{}{"IsPublished": true}
	rns, err := s.repo.FindAll(orgId, page, pageSize, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error finding release notes by organisation ID")
		return nil, err
	}

	for _, rn := range rns.Items {
		if rn.ReleaseDate != nil {
			rd := (*rn.ReleaseDate)[:10]
			rn.ReleaseDate = &rd
			imgUrl, err := s.repo.GetImageUrl(createImgPath(orgId, rn.ID.String()))
			if err != nil {
				log.Error().Err(err).Msg("Error getting image URL")
			}
			rn.ImageUrl = imgUrl
		}
	}

	return rns, nil
}

func (s *service) GetOne(id, orgId string) (*ReleaseNote, error) {
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

	imgUrl, err := s.repo.GetImageUrl(createImgPath(orgId, id))
	if err != nil {
		log.Error().Err(err).Msg("Error getting image URL")
	}
	rn.ImageUrl = imgUrl
	return rn, nil
}

func (s *service) Update(id uuid.UUID, rn *ReleaseNote, imgInput *ImageInput) error {
	log.Trace().Msg("UpdateWithImg")
	// Start a transaction
	tx := s.repo.db.StartTransaction()

	// Update image
	if imgInput != nil {
		path := createImgPath(rn.OrganisationID.String(), id.String())
		if imgInput.ShoudDeleteImage {
			if err := s.repo.DeleteImage(path); err != nil {
				log.Error().Err(err).Msg("Error deleting image")
				tx.Rollback()
				return err
			}
		} else if imgInput.ImgData != nil {
			processedImg, err := imgUtil.DecodeProcessEncode(imgInput.ImgData, &imgProcessConfig)
			if err != nil {
				log.Error().Err(err).Msg("Error processing image")
				tx.Rollback()
				return err
			}
			log.Debug().Str("path", path).Msg("Updating image")
			if err := s.repo.UpdateImage(path, processedImg); err != nil {
				log.Error().Err(err).Msg("Error updating image")
				tx.Rollback()
				return err
			}
		}
	}

	// Update release note
	if err := s.repo.Update(id, rn, tx.Tx); err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		if imgInput != nil && imgInput.ImgData != nil {
			path := createImgPath(rn.OrganisationID.String(), id.String())
			s.repo.DeleteImage(path)
		}
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) ChangePublishedStatus(id uuid.UUID, published bool) error {
	log.Trace().Bool("published", published).Msg("ChangePublishedStatus")
	if err := s.repo.UpdateWithNil(id, map[string]interface{}{"IsPublished": published}); err != nil {
		log.Error().Err(err).Msg("Error updating published status")
		return err
	}
	return nil
}

func (s *service) Delete(id uuid.UUID) error {
	log.Trace().Msg("Delete")
	return s.repo.Delete(id, nil)
}
