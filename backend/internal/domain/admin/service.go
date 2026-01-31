package admin

import (
	"errors"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

type service struct {
	repo        repository
	adminUserID string
}

// NewService creates a new admin service.
func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	cfg := config.New()
	return &service{
		repo:        r,
		adminUserID: cfg.AdminUserID,
	}
}

// IsAdminUser checks if the provided user ID is the admin user
func (s *service) IsAdminUser(userID uuid.UUID) bool {
	log.Trace().Str("userID", userID.String()).Msg("IsAdminUser")
	return IsAdmin(userID, s.adminUserID)
}

// GetAllOrganisations retrieves all organisations if the user is an admin
func (s *service) GetAllOrganisations(userID uuid.UUID) ([]*organisation.Organisation, error) {
	log.Trace().Str("userID", userID.String()).Msg("GetAllOrganisations")

	if !s.IsAdminUser(userID) {
		log.Warn().Str("userID", userID.String()).Msg("Unauthorized access attempt to admin functionality")
		return nil, errors.New("unauthorized access")
	}

	return s.repo.GetAllOrganisations()
}

// GetOrganisationWithUsers retrieves an organisation with its users if the user is an admin
func (s *service) GetOrganisationWithUsers(userID, orgID uuid.UUID) (*organisation.Organisation, []*organisation.OrganisationUser, error) {
	log.Trace().Str("userID", userID.String()).Str("orgID", orgID.String()).Msg("GetOrganisationWithUsers")

	if !s.IsAdminUser(userID) {
		log.Warn().Str("userID", userID.String()).Msg("Unauthorized access attempt to admin functionality")
		return nil, nil, errors.New("unauthorized access")
	}

	return s.repo.GetOrganisationWithUsers(orgID)
}
