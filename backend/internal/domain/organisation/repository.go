package organisation

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

func (r *repository) FindOrg(orgId uuid.UUID) (*Organisation, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("FindOrgById")
	var org Organisation

	if err := r.db.Client.First(&org, "id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation")
		return nil, err
	}
	return &org, nil
}

func (r *repository) FindOrgByExternalId(externalId uuid.UUID) (*Organisation, error) {
	log.Trace().Str("externalId", externalId.String()).Msg("FindOrgByExternalId")
	var org Organisation

	if err := r.db.Client.First(&org, "external_id = ?", externalId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation by external id")
		return nil, err
	}
	return &org, nil
}

func (r *repository) UpdateOrg(orgId uuid.UUID, org *Organisation) error {
	log.Trace().Str("orgId", orgId.String()).Msg("UpdateOrg")
	return r.db.Client.Model(&Organisation{}).Where("id = ?", orgId).Updates(org).Error
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

func (r *repository) FindOrgUser(orgUserId uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("orgUserId", orgUserId.String()).Msg("FindByOrgUserId")
	var ou OrganisationUser
	if err := r.db.Client.First(&ou, orgUserId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation user")
		return nil, err
	}
	return &ou, nil
}

func (r *repository) FindOrgUserByUserId(userId uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("userId", userId.String()).Msg("FindByUserId")
	var ou OrganisationUser

	if err := r.db.Client.Preload("User").First(&ou, "user_id = ?", userId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation user")
		return nil, err
	}
	return &ou, nil
}

func (r *repository) FindOrgUsers(orgId uuid.UUID) ([]*OrganisationUser, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("FindByOrgId")
	var ous []*OrganisationUser

	if err := r.db.Client.Model(&OrganisationUser{}).Preload("User").Find(&ous, "organisation_id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation users by organisation id")
		return nil, err
	}
	log.Debug().Interface("ous", ous).Msg("Organisation users found")
	return ous, nil
}

func (r *repository) DeleteOrgUser(orgUserID uuid.UUID, tx *gorm.DB) error {
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}
	log.Trace().Str("userId", orgUserID.String()).Msg("DeleteByUserId")
	return client.Delete(&OrganisationUser{}, orgUserID).Error
}

func (r *repository) CreateInvite(invite *OrganisationInvite) error {
	log.Trace().Msg("CreateInvite")
	return r.db.Client.Create(invite).Error
}

func (r *repository) FindInvites(orgId uuid.UUID) ([]*OrganisationInvite, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("FindInvitesByOrgId")
	var invites []*OrganisationInvite

	if err := r.db.Client.Model(&OrganisationInvite{}).Find(&invites, "organisation_id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding invites by organisation id")
		return nil, err
	}
	log.Debug().Interface("invites", invites).Msg("Invites found")
	return invites, nil
}

func (r *repository) FindInviteByExternalId(externalId string) (*OrganisationInvite, error) {
	log.Trace().Str("externalId", externalId).Msg("FindInviteByExternalId")
	var invite OrganisationInvite

	if err := r.db.Client.Preload("Organisation").First(&invite, "external_id = ?", externalId).Error; err != nil {
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
