package widgetconfigs

import (
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) Create(cfg *WidgetConfig) (uuid.UUID, error) {
	log.Trace().Msg("Create")
	if err := s.repo.Create(cfg); err != nil {
		log.Error().Err(err).Msg("Error creating widget config")
		return uuid.Nil, err
	}
	return cfg.ID, nil
}

func (s *service) UpdateBaseUrl(orgId string, baseUrl *string) error {
	log.Trace().Msg("UpdateBaseUrl")
	updateMap := map[string]interface{}{
		"ReleasePageBaseUrl": baseUrl,
	}
	if err := s.repo.UpdatePartial(orgId, updateMap); err != nil {
		log.Error().Err(err).Msg("Error updating base URL")
		return err
	}
	return nil
}

func (s *service) Get(orgId string) (*WidgetConfig, error) {
	log.Trace().Str("orgId", orgId).Msg("Get")
	cfg, err := s.repo.Get(orgId)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by organisation ID")
		return nil, err
	}
	return cfg, nil
}

func (s *service) Update(orgId string, cfg *WidgetConfig) error {
	log.Trace().Msg("Update")
	if err := s.repo.Update(orgId, cfg); err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		return err
	}
	return nil
}
