package admin

import (
	"errors"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	"github.com/google/uuid"
)

type service struct {
	repo        repository
	adminUserId string
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	cfg := config.New()
	return &service{
		repo:        r,
		adminUserId: cfg.AdminUserId,
	}
}

// IsAdminUser checks if the provided user ID is the admin user
func (s *service) IsAdminUser(userId uuid.UUID) bool {
	log.Trace().Str("userId", userId.String()).Msg("IsAdminUser")
	return IsAdmin(userId, s.adminUserId)
}

// GetAllOrganisations retrieves all organisations if the user is an admin
func (s *service) GetAllOrganisations(userId uuid.UUID) ([]*organisation.Organisation, error) {
	log.Trace().Str("userId", userId.String()).Msg("GetAllOrganisations")

	if !s.IsAdminUser(userId) {
		log.Warn().Str("userId", userId.String()).Msg("Unauthorized access attempt to admin functionality")
		return nil, errors.New("unauthorized access")
	}

	return s.repo.GetAllOrganisations()
}

// GetOrganisationWithUsers retrieves an organisation with its users if the user is an admin
func (s *service) GetOrganisationWithUsers(userId, orgId uuid.UUID) (*organisation.Organisation, []*organisation.OrganisationUser, error) {
	log.Trace().Str("userId", userId.String()).Str("orgId", orgId.String()).Msg("GetOrganisationWithUsers")

	if !s.IsAdminUser(userId) {
		log.Warn().Str("userId", userId.String()).Msg("Unauthorized access attempt to admin functionality")
		return nil, nil, errors.New("unauthorized access")
	}

	return s.repo.GetOrganisationWithUsers(orgId)
}

// GetSubscriptions retrieves all subscriptions for an organisation
func (s *service) GetSubscriptions(userId, orgId uuid.UUID) ([]*subscription.Subscription, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetSubscriptions")

	if !s.IsAdminUser(userId) {
		log.Warn().Str("userId", userId.String()).Msg("Unauthorized access attempt to admin functionality")
		return nil, errors.New("unauthorized access")
	}

	return s.repo.GetSubscriptions(orgId)
}

// DeleteSubscription deletes a subscription if the user is an admin
func (s *service) DeleteSubscription(userId, subscriptionId uuid.UUID) error {
	log.Trace().Str("userId", userId.String()).Str("subscriptionId", subscriptionId.String()).Msg("DeleteSubscription")

	if !s.IsAdminUser(userId) {
		log.Warn().Str("userId", userId.String()).Msg("Unauthorized access attempt to admin functionality")
		return errors.New("unauthorized access")
	}

	return s.repo.DeleteSubscription(subscriptionId)
}
