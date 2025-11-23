package widgetconfigs

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/google/uuid"
)

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

func (r *repository) Create(cfg *WidgetConfig) error {
	log.Trace().Msg("Create")
	client := r.db.Client
	if err := client.Create(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error saving widget config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("Widget config created")
	return nil
}

func (r *repository) Update(orgId uuid.UUID, cfg *WidgetConfig) error {
	log.Trace().Msg("Update")
	client := r.db.Client
	if err := client.Model(&WidgetConfig{}).Select(
		"Title",
		"Description",
		"WidgetBorderRadius",
		"WidgetBorderColor",
		"WidgetBorderWidth",
		"WidgetBgColor",
		"WidgetTextColor",
		"WidgetType",
		"ReleaseNoteBorderRadius",
		"ReleaseNoteBorderColor",
		"ReleaseNoteBorderWidth",
		"ReleaseNoteBgColor",
		"ReleaseNoteTextColor",
		"ReleaseNoteCtaText",
		"EnableLikes",
		"LikeButtonText",
		"UnlikeButtonText",
	).Where("organisation_id = ?", orgId).Updates(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("Widget config updated")
	return nil
}

func (r *repository) UpdatePartial(orgId uuid.UUID, fields map[string]interface{}) error {
	log.Trace().Msg("UpdateFields")
	client := r.db.Client
	if err := client.Model(&WidgetConfig{}).Where("organisation_id = ?", orgId).Updates(fields).Error; err != nil {
		log.Error().Err(err).Msg("Error updating widget config fields")
		return err
	}
	log.Debug().Interface("fields", fields).Msg("Widget config fields updated")
	return nil
}

func (r *repository) Get(orgId uuid.UUID) (*WidgetConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("Get")
	var cfg WidgetConfig
	if err := r.db.Client.Model(&WidgetConfig{}).Where("organisation_id = ?", orgId).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}
