package releasenotemetrics

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/logger"
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

func (r *repository) Create(metric *ReleaseNoteMetric) error {
	log.Trace().Interface("metric", metric).Msg("Create")
	if err := r.db.Client.Create(metric).Error; err != nil {
		log.Error().Err(err).Msg("Error creating metric")
		return err
	}
	return nil
}

func (r *repository) FindByReleaseNoteID(releaseNoteID uuid.UUID) ([]ReleaseNoteMetric, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("FindByReleaseNoteID")
	var metrics []ReleaseNoteMetric
	if err := r.db.Client.Where("release_note_id = ?", releaseNoteID).Find(&metrics).Error; err != nil {
		log.Error().Err(err).Msg("Error finding metrics")
		return nil, err
	}
	return metrics, nil
}

func (r *repository) FindByOrgID(orgID uuid.UUID) ([]ReleaseNoteMetric, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("FindByOrgID")
	var metrics []ReleaseNoteMetric
	if err := r.db.Client.Where("organisation_id = ?", orgID).Find(&metrics).Error; err != nil {
		log.Error().Err(err).Msg("Error finding metrics")
		return nil, err
	}
	return metrics, nil
}
