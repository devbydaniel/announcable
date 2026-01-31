package widgetconfigs

import (
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

// NewService creates a new widget config service with the given repository.
func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

// DefaultConfig returns a WidgetConfig with default styling values for the given organisation.
func DefaultConfig(orgID uuid.UUID) *WidgetConfig {
	log.Trace().Str("orgID", orgID.String()).Msg("DefaultConfig")
	return &WidgetConfig{
		OrganisationID:          orgID,
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
		EnableLikes:             true,
		LikeButtonText:          "Like",
		UnlikeButtonText:        "Unlike",
	}
}

func (s *service) Init(orgID uuid.UUID) (*WidgetConfig, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Init")
	cfg := DefaultConfig(orgID)
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

func (s *service) UpdateBaseURL(orgID uuid.UUID, baseURL *string) error {
	log.Trace().Msg("UpdateBaseURL")
	updateMap := map[string]interface{}{
		"ReleasePageBaseURL": baseURL,
	}
	if err := s.repo.UpdatePartial(orgID, updateMap); err != nil {
		log.Error().Err(err).Msg("Error updating base URL")
		return err
	}
	return nil
}

func (s *service) Get(orgID uuid.UUID) (*WidgetConfig, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Get")
	cfg, err := s.repo.Get(orgID)
	if err != nil {
		log.Error().Err(err).Msg("Error finding widget config by organisation ID")
		return nil, err
	}
	return cfg, nil
}

func (s *service) Update(orgID uuid.UUID, cfg *WidgetConfig) error {
	log.Trace().Msg("Update")
	if err := s.repo.Update(orgID, cfg); err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		return err
	}
	return nil
}
