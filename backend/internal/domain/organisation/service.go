package organisation

import (
	"errors"
	"regexp"
	"time"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/domain/rbac"
	"github.com/devbydaniel/announcable/internal/domain/user"
	"github.com/devbydaniel/announcable/internal/email"
	"github.com/devbydaniel/announcable/internal/random"
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

// NewService creates a new organisation service.
func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) IsValidOrgName(name string) error {
	log.Trace().Str("name", name).Msg("IsValidOrgName")
	if len(name) <= 0 {
		return errors.New("Organisation name is required")
	}
	if len(name) > 100 {
		return errors.New("Organisation name is too long")
	}
	// only allow alphanumeric characters and spaces
	if !regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString(name) {
		return errors.New("Organisation name can only contain alphanumeric characters and spaces")
	}
	return nil
}

func (s *service) OrgNameExists(name string) bool {
	org, _ := s.repo.FindOrgByName(name)
	return org != nil
}

func (s *service) CreateOrgWithAdmin(name string, user *user.User) (*OrganisationUser, error) {
	log.Trace().Str("name", name).Str("user", user.Email).Msg("CreateOrgWithAdmin")
	org, err := New(name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateOrg(org); err != nil {
		return nil, err
	}
	ou := Connect(org, user, rbac.RoleAdmin)
	log.Debug().Interface("ou", ou).Msg("OrganisationUser")

	if err := s.repo.SaveOrgUser(ou, nil); err != nil {
		return nil, err
	}

	return ou, nil
}

func (s *service) GetOrgUser(orgUserID uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("orgUserID", orgUserID.String()).Msg("GetOrgUser")
	return s.repo.FindOrgUser(orgUserID)
}

func (s *service) GetOrgUserByUserID(userID uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("userID", userID.String()).Msg("GetOrgUser")
	return s.repo.FindOrgUserByUserID(userID)
}

func (s *service) GetOrgUsers(orgID uuid.UUID) ([]*OrganisationUser, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetByOrgID")
	return s.repo.FindOrgUsers(orgID)
}

func (s *service) RemoveFromOrg(orgUserID uuid.UUID) error {
	log.Trace().Str("userID", orgUserID.String()).Msg("DeleteByUserID")
	return s.repo.DeleteOrgUser(orgUserID, nil)
}

func (s *service) InviteUser(orgID uuid.UUID, emailAddr string, role rbac.Role) (string, error) {
	log.Trace().Str("orgID", orgID.String()).Str("email", emailAddr).Msg("InviteUser")
	org, err := s.repo.FindOrg(orgID)
	if err != nil {
		return "", err
	}
	token := random.CreateRandomToken()
	expiredAt := time.Now().Add(time.Hour * 24).UnixMilli()
	invite := OrganisationInvite{
		OrganisationID: orgID,
		Email:          emailAddr,
		Role:           role,
		ExpiresAt:      expiredAt,
		ExternalID:     random.EncodeToken(token),
	}
	if err := s.repo.CreateInvite(&invite); err != nil {
		return "", err
	}

	cfg := config.New()
	inviteAcceptURL := cfg.BaseURL + "/invite-accept/" + token

	// Only send email if enabled
	if cfg.IsEmailEnabled() {
		emailConfig := email.UserInviteConfig{
			To:               emailAddr,
			OrganisationName: org.Name,
			ActionURL:        inviteAcceptURL,
		}
		if err := email.SendUserInvite(&emailConfig); err != nil {
			return "", err
		}
	}

	return inviteAcceptURL, nil
}

func (s *service) GetInvites(orgID uuid.UUID) ([]*OrganisationInvite, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetInvites")
	return s.repo.FindInvites(orgID)
}

func (s *service) GetInviteWithToken(token string) (*OrganisationInvite, error) {
	log.Trace().Str("token", token).Msg("GetInviteWithToken")
	externalID := random.EncodeToken(token)
	return s.repo.FindInviteByExternalID(externalID)
}

func (s *service) DeleteInvite(id uuid.UUID) error {
	log.Trace().Msg("DeleteInvite")
	return s.repo.DeleteInvite(id, nil)
}

func (s *service) AcceptInvite(invite *OrganisationInvite, user *user.User) error {
	log.Trace().Msg("AcceptInvite")
	tx := s.repo.db.StartTransaction()
	ou := Connect(&invite.Organisation, user, invite.Role)
	if err := s.repo.SaveOrgUser(ou, tx.Tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := s.repo.DeleteInvite(invite.ID, tx.Tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) GetExternalID(orgID uuid.UUID) (uuid.UUID, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetExternalOrgID")
	org, err := s.repo.FindOrg(orgID)
	if err != nil {
		return uuid.Nil, err
	}
	return org.ExternalID, nil
}

func (s *service) GetOrg(orgID uuid.UUID) (*Organisation, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetOrg")
	return s.repo.FindOrg(orgID)
}

func (s *service) GetOrgByExternalID(externalID uuid.UUID) (*Organisation, error) {
	log.Trace().Str("externalID", externalID.String()).Msg("GetOrgByExternalID")
	return s.repo.FindOrgByExternalID(externalID)
}

func (s *service) UpdateOrg(orgID uuid.UUID, org *Organisation) error {
	log.Trace().Str("orgID", orgID.String()).Msg("UpdateOrg")
	return s.repo.UpdateOrg(orgID, org)
}

func (s *service) RegenerateExternalID(orgID uuid.UUID) (uuid.UUID, error) {
	log.Trace().Msg("RegenerateExternalID")
	externalID, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("Error generating external ID")
		return uuid.Nil, err
	}
	if err := s.repo.UpdateOrg(orgID, &Organisation{ExternalID: externalID}); err != nil {
		log.Error().Err(err).Msg("Error updating external ID")
		return uuid.Nil, err
	}
	return externalID, nil
}
