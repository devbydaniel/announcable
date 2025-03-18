package releasenotelikes

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/google/uuid"
)

var log = logger.Get()

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

func (r *repository) Create(like *ReleaseNoteLike) error {
	log.Trace().Interface("like", like).Msg("Create")
	if err := r.db.Client.Create(like).Error; err != nil {
		log.Error().Err(err).Msg("Error creating like")
		return err
	}
	return nil
}

func (r *repository) Delete(like *ReleaseNoteLike) error {
	log.Trace().Interface("like", like).Msg("Delete")
	if err := r.db.Client.Delete(like).Error; err != nil {
		log.Error().Err(err).Msg("Error deleting like")
		return err
	}
	return nil
}

func (r *repository) FindByReleaseNoteID(releaseNoteID uuid.UUID) ([]ReleaseNoteLike, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("FindByReleaseNoteID")
	var likes []ReleaseNoteLike
	if err := r.db.Client.Where("release_note_id = ?", releaseNoteID).Find(&likes).Error; err != nil {
		log.Error().Err(err).Msg("Error finding likes")
		return nil, err
	}
	return likes, nil
}

func (r *repository) FindByOrgID(orgID uuid.UUID) ([]ReleaseNoteLike, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("FindByOrgID")
	var likes []ReleaseNoteLike
	if err := r.db.Client.Where("organisation_id = ?", orgID).Find(&likes).Error; err != nil {
		log.Error().Err(err).Msg("Error finding likes")
		return nil, err
	}
	return likes, nil
}

func (r *repository) FindByReleaseNoteAndClientID(releaseNoteID uuid.UUID, clientID string) (*ReleaseNoteLike, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Str("clientID", clientID).Msg("FindByReleaseNoteAndClientID")
	var like ReleaseNoteLike
	if err := r.db.Client.Where("release_note_id = ? AND client_id = ?", releaseNoteID, clientID).First(&like).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding like")
		return nil, err
	}
	return &like, nil
} 