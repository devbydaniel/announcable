package widgetconfigs

import (
	"errors"

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

func (r *repository) Update(orgId string, cfg *WidgetConfig) error {
	log.Trace().Msg("Update")
	client := r.db.Client
	if err := client.Model(&WidgetConfig{}).Select("Title", "Description", "WidgetBorderRadius", "WidgetBorderColor", "WidgetBorderWidth", "WidgetBgColor", "WidgetTextColor", "WidgetType", "ReleaseNoteBorderRadius", "ReleaseNoteBorderColor", "ReleaseNoteBorderWidth", "ReleaseNoteBgColor", "ReleaseNoteTextColor", "ReleaseNoteCtaText").Where("organisation_id = ?", uuid.MustParse(orgId)).Updates(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error updating widget config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("Widget config updated")
	return nil
}

func (r *repository) UpdatePartial(orgId string, fields map[string]interface{}) error {
	log.Trace().Msg("UpdateFields")
	client := r.db.Client
	if err := client.Model(&WidgetConfig{}).Where("organisation_id = ?", uuid.MustParse(orgId)).Updates(fields).Error; err != nil {
		log.Error().Err(err).Msg("Error updating widget config fields")
		return err
	}
	log.Debug().Interface("fields", fields).Msg("Widget config fields updated")
	return nil
}

func (r *repository) Get(orgId string) (*WidgetConfig, error) {
	log.Trace().Str("orgId", orgId).Msg("Get")
	var cfg WidgetConfig
	defaultCfg := WidgetConfig{
		Title:                "Release Notes",
		Description:          "Stay up to date with our latest releases",
		WidgetBorderRadius:   12,
		WidgetBorderColor:    "#000000",
		WidgetBorderWidth:    1,
		WidgetBgColor:        "#f8f9fa",
		WidgetTextColor:      "#000000",
		WidgetType:           "modal",
		ReleaseNoteBgColor:   "#ffffff",
		ReleaseNoteTextColor: "#000000",
		ReleaseNoteCtaText:   "Read more",
		OrganisationID:       uuid.MustParse(orgId),
	}

	if err := r.db.Client.Model(&WidgetConfig{}).Where("organisation_id = ?", uuid.MustParse(orgId)).First(&cfg).Error; err != nil {
		if errors.Is(err, r.db.ErrRecordNotFound) {
			log.Debug().Msg("Widget config not found, creating...")
			r.Create(&defaultCfg)
			return &defaultCfg, nil
		}
		log.Error().Err(err).Msg("Error finding widget config by organisation id")
		return nil, err
	}
	log.Debug().Interface("cfg", cfg).Msg("Widget config found")
	return &cfg, nil
}
