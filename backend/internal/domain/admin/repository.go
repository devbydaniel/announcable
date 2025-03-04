package admin

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	"github.com/google/uuid"
)

type repository struct {
	db *database.DB
}

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
func (r *repository) GetOrganisationWithUsers(orgId uuid.UUID) (*organisation.Organisation, []*organisation.OrganisationUser, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetOrganisationWithUsers")
	var org organisation.Organisation
	var orgUsers []*organisation.OrganisationUser

	if err := r.db.Client.First(&org, "id = ?", orgId.String()).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation")
		return nil, nil, err
	}

	if err := r.db.Client.Preload("User").Find(&orgUsers, "organisation_id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding organisation users")
		return &org, nil, err
	}

	return &org, orgUsers, nil
}

// GetSubscriptions retrieves all subscriptions for an organisation
func (r *repository) GetSubscriptions(orgId uuid.UUID) ([]*subscription.Subscription, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetSubscriptions")
	var subscriptions []*subscription.Subscription

	if err := r.db.Client.Find(&subscriptions, "organisation_id = ?", orgId).Error; err != nil {
		log.Error().Err(err).Msg("Error finding subscriptions")
		return nil, err
	}

	return subscriptions, nil
}

// DeleteSubscription deletes a subscription by ID
func (r *repository) DeleteSubscription(subscriptionId uuid.UUID) error {
	log.Trace().Str("subscriptionId", subscriptionId.String()).Msg("DeleteSubscription")

	if err := r.db.Client.Delete(&subscription.Subscription{}, "id = ?", subscriptionId).Error; err != nil {
		log.Error().Err(err).Msg("Error deleting subscription")
		return err
	}

	return nil
}