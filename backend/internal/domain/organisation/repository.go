package organisation

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *database.DB
}

// NewRepository creates a new organisation repository.
func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

func (r *repository) CreateOrg(org *Organisation) error {
	log.Trace().Msg("CreateOrg")
	return r.db.Client.Create(org).Error
}

func (r *repository) FindOrgByName(name string) (*Organisation, error) {
	log.Trace().Str("name", name).Msg("FindOrgByName")
	var org Organisation
	if err := r.db.Client.First(&org, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *repository) FindOrg(orgID uuid.UUID) (*Organisation, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("FindOrgById")
	var org Organisation

	if err := r.db.Client.First(&org, "id = ?", orgID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation")
		return nil, err
	}
	return &org, nil
}

func (r *repository) FindOrgByExternalID(externalID uuid.UUID) (*Organisation, error) {
	log.Trace().Str("externalID", externalID.String()).Msg("FindOrgByExternalID")
	var org Organisation

	if err := r.db.Client.First(&org, "external_id = ?", externalID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation by external id")
		return nil, err
	}
	return &org, nil
}

func (r *repository) UpdateOrg(orgID uuid.UUID, org *Organisation) error {
	log.Trace().Str("orgID", orgID.String()).Msg("UpdateOrg")
	return r.db.Client.Model(&Organisation{}).Where("id = ?", orgID).Updates(org).Error
}

func (r *repository) SaveOrgUser(ou *OrganisationUser, tx *gorm.DB) error {
	log.Trace().Msg("Save")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	return client.Create(ou).Error
}

func (r *repository) FindOrgUser(orgUserID uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("orgUserID", orgUserID.String()).Msg("FindByOrgUserId")
	var ou OrganisationUser
	if err := r.db.Client.First(&ou, orgUserID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation user")
		return nil, err
	}
	return &ou, nil
}

func (r *repository) FindOrgUserByUserID(userID uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("userID", userID.String()).Msg("FindByUserId")
	var ou OrganisationUser

	if err := r.db.Client.Preload("User").Preload("Organisation").First(&ou, "user_id = ?", userID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation user")
		return nil, err
	}
	return &ou, nil
}

func (r *repository) FindOrgUsers(orgID uuid.UUID) ([]*OrganisationUser, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("FindByOrgID")
	var ous []*OrganisationUser

	if err := r.db.Client.Model(&OrganisationUser{}).Preload("User").Find(&ous, "organisation_id = ?", orgID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation users by organisation id")
		return nil, err
	}
	log.Debug().Msg("Organisation users found")
	return ous, nil
}

func (r *repository) DeleteOrgUser(orgUserID uuid.UUID, tx *gorm.DB) error {
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	log.Trace().Str("userID", orgUserID.String()).Msg("DeleteByUserID")
	return client.Delete(&OrganisationUser{}, orgUserID).Error
}

func (r *repository) CreateInvite(invite *OrganisationInvite) error {
	log.Trace().Msg("CreateInvite")
	return r.db.Client.Create(invite).Error
}

func (r *repository) FindInvites(orgID uuid.UUID) ([]*OrganisationInvite, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("FindInvitesByOrgID")
	var invites []*OrganisationInvite

	if err := r.db.Client.Model(&OrganisationInvite{}).Find(&invites, "organisation_id = ?", orgID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding invites by organisation id")
		return nil, err
	}
	log.Debug().Interface("invites", invites).Msg("Invites found")
	return invites, nil
}

func (r *repository) FindInviteByExternalID(externalID string) (*OrganisationInvite, error) {
	log.Trace().Str("externalID", externalID).Msg("FindInviteByExternalID")
	var invite OrganisationInvite

	if err := r.db.Client.Preload("Organisation").First(&invite, "external_id = ?", externalID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding invite by external id")
		return nil, err
	}
	return &invite, nil
}

func (r *repository) DeleteInvite(id uuid.UUID, tx *gorm.DB) error {
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	log.Trace().Msg("DeleteInvite")
	return client.Delete(&OrganisationInvite{}, id).Error
}
