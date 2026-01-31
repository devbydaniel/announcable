package releasepageconfig

import (
	"io"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/objstore"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db       *database.DB
	objStore *objstore.ObjStore
	bucket   string
}

// NewRepository creates a new release page config repository with database and object storage access.
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

func (r *repository) Update(orgID uuid.UUID, cfg *ReleasePageConfig, tx *gorm.DB) error {
	log.Trace().Msg("Update")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&ReleasePageConfig{}).Where("organisation_id = ?", orgID).Select("Title", "Description", "BgColor", "TextColor", "TextColorMuted", "BrandPosition", "BackLinkLabel", "BackLinkURL").Updates(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("landing page config updated")
	return nil
}

func (r *repository) UpdateWithNil(orgID uuid.UUID, fields map[string]interface{}, tx *gorm.DB) error {
	log.Trace().Msg("UpdateWithNil")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&ReleasePageConfig{}).Where("organisation_id = ?", orgID).Updates(fields).Error; err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	log.Debug().Interface("fields", fields).Msg("landing page config updated")
	return nil
}

func (r *repository) Get(orgID uuid.UUID) (*ReleasePageConfig, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Get")
	var cfg ReleasePageConfig
	if err := r.db.Client.Model(&ReleasePageConfig{}).Preload("Organisation").Where("organisation_id = ?", orgID).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *repository) GetBySlug(slug string) (*ReleasePageConfig, error) {
	log.Trace().Str("slug", slug).Msg("GetBySlug")
	var cfg ReleasePageConfig
	if err := r.db.Client.Model(&ReleasePageConfig{}).Where("slug = ?", slug).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *repository) GetImageURL(path string) (string, error) {
	log.Trace().Msg("GetImageURL")
	return r.objStore.GetImageURL(r.bucket, path)
}

func (r *repository) UpdateImage(path string, img *io.Reader) error {
	log.Trace().Msg("UpdateImage")
	return r.objStore.UpdateImage(r.bucket, path, img)
}

func (r *repository) DeleteImage(orgID uuid.UUID) error {
	log.Trace().Msg("DeleteImage")
	cfg, err := r.Get(orgID)
	if err != nil {
		log.Error().Err(err).Msg("Error finding landing page config")
		return err
	}
	if err := r.objStore.DeleteImage(r.bucket, cfg.ImagePath); err != nil {
		log.Error().Err(err).Msg("Error deleting image")
		return err
	}
	if err := r.UpdateWithNil(orgID, map[string]interface{}{"ImagePath": nil}, nil); err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	return nil
}
