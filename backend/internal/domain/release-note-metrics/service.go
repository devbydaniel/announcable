package releasenotemetrics

import (
	"github.com/google/uuid"
)

type service struct {
	repo *repository
}

// NewService creates a new release note metrics service.
func NewService(r *repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) CreateMetric(releaseNoteID uuid.UUID, orgID uuid.UUID, clientID string, metricType MetricType) error {
	log.Trace().Msg("CreateMetric")
	metric := &ReleaseNoteMetric{
		ReleaseNoteID:  releaseNoteID,
		OrganisationID: orgID,
		ClientID:       clientID,
		MetricType:     metricType,
	}
	return s.repo.Create(metric)
}

func (s *service) GetViewCount(releaseNoteID uuid.UUID) (int, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("GetViewCount")
	metrics, err := s.repo.FindByReleaseNoteID(releaseNoteID)
	if err != nil {
		return 0, err
	}

	viewCount := 0
	for _, metric := range metrics {
		if metric.MetricType == MetricTypeView {
			viewCount++
		}
	}
	return viewCount, nil
}

func (s *service) GetCtaClickCount(releaseNoteID uuid.UUID) (int, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("GetCtaClickCount")
	metrics, err := s.repo.FindByReleaseNoteID(releaseNoteID)
	if err != nil {
		return 0, err
	}

	clickCount := 0
	for _, metric := range metrics {
		if metric.MetricType == MetricTypeCtaClick {
			clickCount++
		}
	}
	return clickCount, nil
}

func (s *service) GetMetricsByReleaseNote(releaseNoteID uuid.UUID) ([]ReleaseNoteMetric, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("GetMetricsByReleaseNote")
	return s.repo.FindByReleaseNoteID(releaseNoteID)
}

func (s *service) GetMetricsByOrg(orgID uuid.UUID) ([]ReleaseNoteMetric, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetMetricsByOrg")
	return s.repo.FindByOrgID(orgID)
}
