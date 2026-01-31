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
	"github.com/devbydaniel/announcable/internal/util"
	"github.com/google/uuid"
)

type service struct {
	repo repository
}

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

func (s *service) GetOrgUser(orgUserId uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("orgUserId", orgUserId.String()).Msg("GetOrgUser")
	return s.repo.FindOrgUser(orgUserId)
}

func (s *service) GetOrgUserByUserId(userId uuid.UUID) (*OrganisationUser, error) {
	log.Trace().Str("userId", userId.String()).Msg("GetOrgUser")
	return s.repo.FindOrgUserByUserId(userId)
}

func (s *service) GetOrgUsers(orgId uuid.UUID) ([]*OrganisationUser, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetByOrgId")
	return s.repo.FindOrgUsers(orgId)
}

func (s *service) RemoveFromOrg(orgUserId uuid.UUID) error {
	log.Trace().Str("userId", orgUserId.String()).Msg("DeleteByUserId")
	return s.repo.DeleteOrgUser(orgUserId, nil)
}

func (s *service) InviteUser(orgId uuid.UUID, emailAddr string, role rbac.Role) (string, error) {
	log.Trace().Str("orgId", orgId.String()).Str("email", emailAddr).Msg("InviteUser")
	org, err := s.repo.FindOrg(orgId)
	if err != nil {
		return "", err
	}
	token := random.CreateRandomToken()
	expiredAt := time.Now().Add(time.Hour * 24).UnixMilli()
	invite := OrganisationInvite{
		OrganisationID: orgId,
		Email:          emailAddr,
		Role:           role,
		ExpiresAt:      expiredAt,
		ExternalID:     random.EncodeToken(token),
	}
	if err := s.repo.CreateInvite(&invite); err != nil {
		return "", err
	}

	cfg := config.New()
	inviteAcceptUrl := util.BuildURL(cfg.BaseURL, "invite-accept", token)

	// Only send email if enabled
	if cfg.IsEmailEnabled() {
		emailConfig := email.UserInviteConfig{
			To:               emailAddr,
			OrganisationName: org.Name,
			ActionURL:        inviteAcceptUrl,
		}
		if err := email.SendUserInvite(&emailConfig); err != nil {
			return "", err
		}
	}

	return inviteAcceptUrl, nil
}

func (s *service) GetInvites(orgId uuid.UUID) ([]*OrganisationInvite, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetInvites")
	return s.repo.FindInvites(orgId)
}

func (s *service) GetInviteWithToken(token string) (*OrganisationInvite, error) {
	log.Trace().Str("token", token).Msg("GetInviteWithToken")
	externalId := random.EncodeToken(token)
	return s.repo.FindInviteByExternalId(externalId)
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

func (s *service) GetExternalId(orgId uuid.UUID) (uuid.UUID, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetExternalOrgId")
	org, err := s.repo.FindOrg(orgId)
	if err != nil {
		return uuid.Nil, err
	}
	return org.ExternalID, nil
}

func (s *service) GetOrg(orgId uuid.UUID) (*Organisation, error) {
	log.Trace().Str("orgId", orgId.String()).Msg("GetOrg")
	return s.repo.FindOrg(orgId)
}

func (s *service) GetOrgByExternalId(externalId uuid.UUID) (*Organisation, error) {
	log.Trace().Str("externalId", externalId.String()).Msg("GetOrgByExternalId")
	return s.repo.FindOrgByExternalId(externalId)
}

func (s *service) UpdateOrg(orgId uuid.UUID, org *Organisation) error {
	log.Trace().Str("orgId", orgId.String()).Msg("UpdateOrg")
	return s.repo.UpdateOrg(orgId, org)
}

func (s *service) RegenerateExternalId(orgId uuid.UUID) (uuid.UUID, error) {
	log.Trace().Msg("RegenerateExternalId")
	externalId, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("Error generating external ID")
		return uuid.Nil, err
	}
	if err := s.repo.UpdateOrg(orgId, &Organisation{ExternalID: externalId}); err != nil {
		log.Error().Err(err).Msg("Error updating external ID")
		return uuid.Nil, err
	}
	return externalId, nil
}
