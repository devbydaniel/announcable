package user

import (
	"time"

	"github.com/devbydaniel/release-notes-go/config"
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

	if err := r.db.Client.Preload("TosConfirms").Preload("PrivacyPolicyConfirms").First(&u, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Msg("Error finding user")
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	log.Trace().Str("email", email).Msg("FindByEmail")
	var u User

	if err := r.db.Client.First(&u, "email = ?", email).Error; err != nil {
		log.Error().Err(err).Msg("Error finding user")
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

func (r *repository) ConfirmTosNow(id uuid.UUID) (string, error) {
	log.Trace().Str("id", id.String()).Msg("ConfirmTos")
	currentVersion := config.New().Legal.ToSVersion
	tosConfim := &TosConfirm{UserID: id, Version: currentVersion, ConfirmedAt: time.Now()}
	if err := r.db.Client.Create(tosConfim).Error; err != nil {
		log.Error().Err(err).Msg("Error confirming ToS")
		return "", err
	}
	return currentVersion, nil
}

func (r *repository) GetLatestTosVersion(id uuid.UUID) (string, error) {
	log.Trace().Msg("GetLatestTosVersion")
	var tos TosConfirm
	if err := r.db.Client.Model(&TosConfirm{}).Where("user_id = ?", id).Order("confirmed_at desc").First(&tos).Error; err != nil {
		log.Error().Err(err).Msg("Error getting latest ToS version")
		return "", err
	}
	return tos.Version, nil
}

func (r *repository) ConfirmPrivacyPolicyNow(id uuid.UUID) (string, error) {
	log.Trace().Str("id", id.String()).Msg("PrivacyPolicyConfirm")
	currentVersion := config.New().Legal.PPVersion
	ppConfirm := &PrivacyPolicyConfirm{UserID: id, Version: currentVersion, ConfirmedAt: time.Now()}
	if err := r.db.Client.Create(ppConfirm).Error; err != nil {
		log.Error().Err(err).Msg("Error confirming Privacy Policy")
		return "", err
	}
	return currentVersion, nil
}

func (r *repository) GetLatestPrivacyPolicyVersion(id uuid.UUID) (string, error) {
	log.Trace().Msg("GetLatestPrivacyPolicyVersion")
	var pp PrivacyPolicyConfirm
	if err := r.db.Client.Model(&PrivacyPolicyConfirm{}).Where("user_id = ?", id).Order("confirmed_at desc").First(&pp).Error; err != nil {
		log.Error().Err(err).Msg("Error getting latest Privacy Policy version")
		return "", err
	}
	return pp.Version, nil
}
