package admin

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

type repository struct {
	db *database.DB
}

// NewRepository creates a new admin repository.
func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

// GetAllOrganisations retrieves all organisations from the database
func (r *repository) GetAllOrganisations() ([]*organisation.Organisation, error) {
	log.Trace().Msg("GetAllOrganisations")
	var orgs []*organisation.Organisation

	if err := r.db.Client.Find(&orgs).Error; err != nil {
		log.Error().Err(err).Msg("Error finding all organisations")
		return nil, err
	}

	log.Debug().Int("count", len(orgs)).Msg("Organisations found")
	return orgs, nil
}

// GetOrganisationWithUsers retrieves an organisation with its users
func (r *repository) GetOrganisationWithUsers(orgID uuid.UUID) (*organisation.Organisation, []*organisation.OrganisationUser, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetOrganisationWithUsers")
	var org organisation.Organisation
	var orgUsers []*organisation.OrganisationUser

	if err := r.db.Client.First(&org, "id = ?", orgID.String()).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation")
		return nil, nil, err
	}

	if err := r.db.Client.Preload("User").Find(&orgUsers, "organisation_id = ?", orgID).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation users")
		return &org, nil, err
	}

	return &org, orgUsers, nil
}
