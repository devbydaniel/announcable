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

func DefaultConfig(orgId uuid.UUID) *WidgetConfig {
	log.Trace().Str("orgId", orgId.String()).Msg("DefaultConfig")
	return &WidgetConfig{
		OrganisationID:          orgId,
		Title:                   "Release Notes",
		Description:             "See what's new",
		WidgetBorderRadius:      12,
		WidgetBorderColor:       "#7c7f93",
		WidgetBorderWidth:       1,
		WidgetBgColor:           "#dce0e8",
		WidgetTextColor:         "#4c4f69",
		WidgetType:              "modal",
		ReleaseNoteBorderRadius: 12,
		ReleaseNoteBorderColor:  "#7c7f93",
		ReleaseNoteBorderWidth:  1,
		ReleaseNoteBgColor:      "#eff1f5",
		ReleaseNoteTextColor:    "#4c4f69",
		ReleaseNoteCtaText:      "View Release Notes",
	}
}

func (s *service) Init(orgId uuid.UUID) (*WidgetConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("Init")
	cfg := DefaultConfig(orgId)
	if err := s.repo.Create(cfg); err != nil {
		log.Error().Err(err).Msg("Error creating default widget config")
		return cfg, err
	}
	return cfg, nil
}

func (s *service) Create(cfg *WidgetConfig) (uuid.UUID, error) {
	log.Trace().Msg("Create")
	if err := s.repo.Create(cfg); err != nil {
		log.Error().Err(err).Msg("Error creating widget config")
		return uuid.Nil, err
	}
	return cfg.ID, nil
}

func (s *service) UpdateBaseUrl(orgId uuid.UUID, baseUrl *string) error {
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

func (s *service) Get(orgId uuid.UUID) (*WidgetConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("Get")
	cfg, err := s.repo.Get(orgId)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by organisation ID")
		return nil, err
	}
	return cfg, nil
}

func (s *service) Update(orgId uuid.UUID, cfg *WidgetConfig) error {
	log.Trace().Msg("Update")
	if err := s.repo.Update(orgId, cfg); err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		return err
	}
	return nil
}
