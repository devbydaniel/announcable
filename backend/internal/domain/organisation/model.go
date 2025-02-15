package organisation

import (
	"errors"
	"strings"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/rbac"
	"github.com/devbydaniel/release-notes-go/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organisation struct {
	database.BaseModel `gorm:"embedded"`
	Name               string
	ExternalID         uuid.UUID `gorm:"type:uuid;uniqueIndex"`
}

type OrganisationUser struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID
	Organisation       Organisation
	UserID             uuid.UUID
	User               user.User
	Role               rbac.Role `gorm:"type:varchar(255)"`
}

type OrganisationInvite struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID
	Organisation       Organisation
	Email              string    `gorm:"type:varchar(255)"`
	ExternalID         string    `gorm:"type:varchar(255)"`
	Role               rbac.Role `gorm:"type:varchar(255)"`
	ExpiresAt          int64     `gorm:"type:bigint"`
}

func (o *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	log.Trace().Msg("BeforeCreate")
	o.ExternalID = uuid.New()
	return
}

func New(name string) (*Organisation, error) {
	log.Trace().Str("name", name).Msg("New")
	if strings.TrimSpace(name) == "" || len(name) < 3 {
		return nil, errors.New("Please provide an organisation name with at least 3 characters")
	}
	return &Organisation{Name: name}, nil
}

func Connect(org *Organisation, user *user.User, role rbac.Role) *OrganisationUser {
	log.Trace().Str("org", org.Name).Str("user", user.Email).Str("role", role.String()).Msg("Connect")
	return &OrganisationUser{Organisation: *org, User: *user, Role: role}
}
