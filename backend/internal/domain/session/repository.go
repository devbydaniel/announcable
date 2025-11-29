package session

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/google/uuid"
)

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

func (r *repository) Save(s *Session) error {
	log.Trace().Msg("Save")
	if err := r.db.Client.Save(s).Error; err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	log.Debug().Str("ID", s.ID.String()).Msg("Session saved")
	return nil
}

func (r *repository) FindByExternalId(sessionId string) (*Session, error) {
	log.Trace().Str("sessionId", sessionId).Msg("FindBySessionId")
	s := Session{}
	if err := r.db.Client.Where("external_id = ?", sessionId).First(&s).Error; err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	log.Debug().Str("external_id", s.ExternalID).Str("user_id", s.UserID.String()).Str("ID", s.ID.String()).Msg("Found session")
	return &s, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	log.Trace().Str("sessionId", id.String()).Msg("Delete")
	if err := r.db.Client.Where("id = ?", id).Delete(&Session{}).Error; err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}

func (r *repository) DeleteByUserId(userId uuid.UUID) error {
	log.Trace().Str("userId", userId.String()).Msg("DeleteByUserId")
	if err := r.db.Client.Where("user_id = ?", userId).Delete(&Session{}).Error; err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}
