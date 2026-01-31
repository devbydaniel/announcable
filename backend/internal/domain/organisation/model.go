package organisation

import (
	"errors"
	"strings"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/rbac"
	"github.com/devbydaniel/announcable/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organisation represents a tenant organisation.
type Organisation struct {
	database.BaseModel `gorm:"embedded"`
	Name               string
	ExternalID         uuid.UUID `gorm:"type:uuid;uniqueIndex"`
}

// OrganisationUser represents a user's membership in an organisation.
type OrganisationUser struct { //nolint:revive // stutter is acceptable here
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID
	Organisation       Organisation
	UserID             uuid.UUID
	User               user.User
	Role               rbac.Role `gorm:"type:varchar(255)"`
}

// OrganisationInvite represents a pending invitation to join an organisation.
type OrganisationInvite struct { //nolint:revive // stutter is acceptable here
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID
	Organisation       Organisation
	Email              string    `gorm:"type:varchar(255)"`
	ExternalID         string    `gorm:"type:varchar(255)"`
	Role               rbac.Role `gorm:"type:varchar(255)"`
	ExpiresAt          int64     `gorm:"type:bigint"`
}

// BeforeCreate is a GORM hook that generates an external ID before inserting.
func (o *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	log.Trace().Msg("BeforeCreate")
	o.ExternalID = uuid.New()
	return
}

// New creates a new Organisation with the given name.
func New(name string) (*Organisation, error) {
	log.Trace().Str("name", name).Msg("New")
	if strings.TrimSpace(name) == "" || len(name) < 3 {
		return nil, errors.New("Please provide an organisation name with at least 3 characters")
	}
	return &Organisation{Name: name}, nil
}

// Connect creates an OrganisationUser linking a user to an organisation with the given role.
func Connect(org *Organisation, user *user.User, role rbac.Role) *OrganisationUser {
	log.Trace().Str("org", org.Name).Str("user", user.Email).Str("role", role.String()).Msg("Connect")
	return &OrganisationUser{Organisation: *org, User: *user, Role: role}
}
