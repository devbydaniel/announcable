package user

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

func (r *repository) Create(u *User) error {
	log.Trace().Msg("Save")
	if err := r.db.Client.Create(u).Error; err != nil {
		log.Error().Err(err).Msg("")
		r.db.Client.Rollback()
		return err
	}
	return nil
}

func (r *repository) Update(id uuid.UUID, u *User) error {
	log.Trace().Msg("Update")
	if err := r.db.Client.Model(&User{}).Where("id = ?", id).Updates(u).Error; err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}

func (r *repository) UpdateWithNil(id uuid.UUID, data map[string]interface{}) error {
	log.Trace().Msg("UpdateWithNil")
	if err := r.db.Client.Model(&User{}).Where("id = ?", id).Updates(data).Error; err != nil {
		log.Error().Err(err).Msg("Error updating user")
		return err
	}
	return nil
}

func (r *repository) FindById(id uuid.UUID) (*User, error) {
	log.Trace().Str("id", id.String()).Msg("FindById")
	var u User

	if err := r.db.Client.First(&u, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Msg("Error finding user")
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	log.Trace().Str("email", email).Msg("FindByEmail")
	var u User

	if err := r.db.Client.First(&u, "email = ?", email).Error; err != nil {
		log.Warn().Err(err).Msg("Error finding user")
		return nil, err
	}
	return &u, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	log.Trace().Str("id", id.String()).Msg("Delete")
	if err := r.db.Client.Delete(&User{}, id).Error; err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}
