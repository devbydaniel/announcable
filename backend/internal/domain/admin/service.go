package admin

import (
	"errors"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
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
