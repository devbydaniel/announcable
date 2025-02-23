package releasepageconfig

import (
	"errors"
	"io"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/objstore"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db       *database.DB
	objStore *objstore.ObjStore
	bucket   string
}

func NewRepository(db *database.DB, objStore *objstore.ObjStore) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db, objStore: objStore, bucket: objstore.LandingPageBucket.String()}
}

func (r *repository) Create(cfg *ReleasePageConfig, tx *gorm.DB) error {
	log.Trace().Msg("Create")

	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Create(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error saving landing page config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("landing page config created")
	return nil
}

func (r *repository) Update(orgId uuid.UUID, cfg *ReleasePageConfig, tx *gorm.DB) error {
	log.Trace().Msg("Update")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&ReleasePageConfig{}).Where("organisation_id = ?", orgId).Select("Title", "Description", "BgColor", "TextColor", "TextColorMuted", "BrandPosition", "BackLinkLabel", "BackLinkUrl").Updates(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("landing page config updated")
	return nil
}

func (r *repository) UpdateWithNil(orgId uuid.UUID, fields map[string]interface{}, tx *gorm.DB) error {
	log.Trace().Msg("UpdateWithNil")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&ReleasePageConfig{}).Where("organisation_id = ?", orgId).Updates(fields).Error; err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	log.Debug().Interface("fields", fields).Msg("landing page config updated")
	return nil
}

func (r *repository) Get(orgId uuid.UUID) (*ReleasePageConfig, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("Get")
	var cfg ReleasePageConfig
	defaultCfg := ReleasePageConfig{
		Title:          "Release Notes",
		Description:    "Stay up to date with our latest releases",
		BgColor:        "#f8f9fa",
		TextColor:      "#000000",
		TextColorMuted: "#6c757d",
		BrandPosition:  string(BrandPositionTop),
		OrganisationID: orgId,
	}

	if err := r.db.Client.Model(&ReleasePageConfig{}).Where("organisation_id = ?", orgId).First(&cfg).Error; err != nil {
		if errors.Is(err, r.db.ErrRecordNotFound) {
			log.Debug().Msg("landing page config not found, creating...")
			r.Create(&defaultCfg, nil)
			return &defaultCfg, nil
		}
		log.Error().Err(err).Msg("Error finding landing page config by organisation id")
		return nil, err
	}
	return &cfg, nil
}

func (r *repository) GetImageUrl(path string) (string, error) {
	log.Trace().Msg("GetImageUrl")
	return r.objStore.GetImageUrl(r.bucket, path)
}

func (r *repository) UpdateImage(path string, img *io.Reader) error {
	log.Trace().Msg("UpdateImage")
	return r.objStore.UpdateImage(r.bucket, path, img)
}

func (r *repository) DeleteImage(orgId uuid.UUID) error {
	log.Trace().Msg("DeleteImage")
	cfg, err := r.Get(orgId)
	if err != nil {
		log.Error().Err(err).Msg("Error finding landing page config")
		return err
	}
	if err := r.objStore.DeleteImage(r.bucket, cfg.ImagePath); err != nil {
		log.Error().Err(err).Msg("Error deleting image")
		return err
	}
	if err := r.UpdateWithNil(orgId, map[string]interface{}{"ImagePath": nil}, nil); err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	return nil
}
