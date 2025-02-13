package releasenotes

import (
	"io"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/objstore"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db       *database.DB
	tx       *database.Transaction
	objStore *objstore.ObjStore
	bucket   string
}

func (r *repository) StartTransaction() {
	r.tx = r.db.StartTransaction()
}

func (r *repository) Commit() {
	r.tx.Commit()
}

func (r *repository) Rollback() {
	r.tx.Rollback()
}

func NewRepository(db *database.DB, objStore *objstore.ObjStore) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db, objStore: objStore, bucket: objstore.ReleaseNotesBucket.String()}
}

func (r *repository) Create(rn *ReleaseNote, tx *gorm.DB) (uuid.UUID, error) {
	log.Trace().Msg("Create")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	if err := client.Create(rn).Error; err != nil {
		log.Error().Err(err).Msg("Error saving release note")
		return uuid.Nil, err
	}
	log.Debug().Interface("rn", rn).Msg("Release note created")
	return rn.ID, nil
}

func (r *repository) Update(id uuid.UUID, rn *ReleaseNote, tx *gorm.DB) error {
	log.Trace().Msg("Update")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	if err := client.Model(&ReleaseNote{}).Where("id = ?", id).Select("Title", "DescriptionShort", "DescriptionLong", "ReleaseDate", "CtaLabelOverride", "CtaUrlOverride", "HideCta", "AttentionMechanism").Updates(rn).Error; err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		return err
	}
	log.Debug().Interface("rn", rn).Msg("Release note updated")
	return nil
}

func (r *repository) UpdateWithNil(id uuid.UUID, data map[string]interface{}) error {
	log.Trace().Msg("UpdatePartial")
	client := r.db.Client
	if err := client.Model(&ReleaseNote{}).Where("id = ?", id).Updates(data).Error; err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		return err
	}
	log.Debug().Msg("Release note updated")
	return nil
}

func (r *repository) FindAll(orgId string) ([]*ReleaseNote, error) {
	log.Trace().Str("orgId", orgId).Msg("FindByOrganisationId")
	var rns []*ReleaseNote

	if err := r.db.Client.Find(&rns, "organisation_id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding release notes by organisation id")
		return nil, err
	}
	return rns, nil
}

func (r *repository) FindOne(id uuid.UUID) (*ReleaseNote, error) {
	log.Trace().Msg("FindById")
	rn := &ReleaseNote{}
	if err := r.db.Client.First(rn, id).Error; err != nil {
		log.Error().Err(err).Msg("Error finding release note by id")
		return nil, err
	}
	return rn, nil
}

func (r *repository) Delete(id uuid.UUID, tx *gorm.DB) error {
	log.Trace().Msg("Delete")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	if err := client.Delete(&ReleaseNote{}, id).Error; err != nil {
		log.Error().Err(err).Msg("Error deleting release note")
		return err
	}
	return nil
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
