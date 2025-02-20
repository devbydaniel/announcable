package releasenotes

import (
	"io"
	"math"

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

func (r *repository) UpdateWithNil(id uuid.UUID, data map[string]interface{}, tx *gorm.DB) error {
	log.Trace().Msg("UpdatePartial")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	if err := client.Model(&ReleaseNote{}).Where("id = ?", id).Updates(data).Error; err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		return err
	}
	log.Debug().Msg("Release note updated")
	return nil
}

func (r *repository) GetStatus(orgId string) ([]*ReleaseNoteStatus, error) {
	log.Trace().Str("orgId", orgId).Msg("GetStatus")
	var statuses []*ReleaseNoteStatus
	if err := r.db.Client.Model(&ReleaseNote{}).Select("updated_at, attention_mechanism").Where("organisation_id = ?", orgId).Find(&statuses).Error; err != nil {
		log.Error().Err(err).Msg("Error getting release note statuses")
		return nil, err
	}
	return statuses, nil
}

func (r *repository) FindAll(orgId string, page, pageSize int, filters map[string]interface{}) (*PaginatedReleaseNotes, error) {
	log.Trace().Str("orgId", orgId).Int("page", page).Int("pageSize", pageSize).Msg("FindByOrganisationId")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Base query conditions
	query := r.db.Client.Model(&ReleaseNote{}).Where("organisation_id = ?", orgId)

	if len(filters) > 0 {
		// Apply additional filters
		for key, value := range filters {
			query = query.Where(key+" = ?", value)
		}
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		log.Error().Err(err).Msg("Error counting release notes")
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	var rns []*ReleaseNote
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("release_date desc").Find(&rns).Error; err != nil {
		log.Error().Err(err).Msg("Error finding release notes by organisation id")
		return nil, err
	}

	return &PaginatedReleaseNotes{
		Items:      rns,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}, nil
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

func (r *repository) UpdateImage(id uuid.UUID, img *io.Reader, path string, tx *gorm.DB) error {
	log.Trace().Msg("UpdateImage")
	if err := r.UpdateWithNil(id, map[string]interface{}{"ImagePath": path}, tx); err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		return err
	}
	if err := r.objStore.UpdateImage(r.bucket, path, img); err != nil {
		log.Error().Err(err).Msg("Error updating image")
		return err
	}
	return nil
}

func (r *repository) DeleteImage(id uuid.UUID, tx *gorm.DB) error {
	log.Trace().Msg("DeleteImage")
	rn, err := r.FindOne(id)
	if err != nil {
		log.Error().Err(err).Msg("Error finding release note")
		return err
	}
	if err := r.objStore.DeleteImage(r.bucket, rn.ImagePath); err != nil {
		log.Error().Err(err).Msg("Error deleting image")
		return err
	}
	if err := r.UpdateWithNil(id, map[string]interface{}{"ImagePath": nil}, tx); err != nil {
		log.Error().Err(err).Msg("Error updating release note")
		return err
	}
	return nil
}
