package lpconfigs

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

func (r *repository) Create(cfg *LpConfig, tx *gorm.DB) error {
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

func (r *repository) Update(orgId string, cfg *LpConfig, tx *gorm.DB) error {
	log.Trace().Msg("Update")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&LpConfig{}).Where("organisation_id = ?", uuid.MustParse(orgId)).Updates(cfg).Error; err != nil {
		log.Error().Err(err).Msg("Error updating landing page config")
		return err
	}
	log.Debug().Interface("cfg", cfg).Msg("landing page config updated")
	return nil
}

func (r *repository) Get(orgId string) (*LpConfig, error) {
	log.Trace().Str("orgId", orgId).Msg("Get")
	var cfg LpConfig
	defaultCfg := LpConfig{
		Title:          "Release Notes",
		Description:    "Stay up to date with our latest releases",
		BgColor:        "#f8f9fa",
		TextColor:      "#000000",
		TextColorMuted: "#6c757d",
		BrandPosition:  string(BrandPositionTop),
		OrganisationID: uuid.MustParse(orgId),
	}

	if err := r.db.Client.Model(&LpConfig{}).Where("organisation_id = ?", uuid.MustParse(orgId)).First(&cfg).Error; err != nil {
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

func (r *repository) DeleteImage(path string) error {
	log.Trace().Msg("DeleteImage")
	return r.objStore.DeleteImage(r.bucket, path)
}
