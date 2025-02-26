package user

import (
	"errors"
	"regexp"
	"time"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/google/uuid"
)

type User struct {
	database.BaseModel    `gorm:"embedded"`
	Email                 string `gorm:"unique"`
	Password              string
	EmailVerified         bool                   `gorm:"default:false"`
	TosConfirms           []TosConfirm           `gorm:"foreignKey:UserID"`
	PrivacyPolicyConfirms []PrivacyPolicyConfirm `gorm:"foreignKey:UserID"`
}

type TosConfirm struct {
	database.BaseModel `gorm:"embedded"`
	UserID             uuid.UUID
	User               User
	Version            string    `gorm:"type:varchar(255)"`
	ConfirmedAt        time.Time `gorm:"type:timestamptz"`
}

type PrivacyPolicyConfirm struct {
	database.BaseModel `gorm:"embedded"`
	UserID             uuid.UUID
	User               User
	Version            string    `gorm:"type:varchar(255)"`
	ConfirmedAt        time.Time `gorm:"type:timestamptz"`
}

func New(email, password string) (*User, error) {
	log.Trace().Str("email", email).Msg("New")
	emailRegex, err := regexp.Compile(`^.+@.+\.[A-Za-z]{2,}$`)
	if err != nil {
		log.Error().Err(err).Msg("Failed to compile email regex")
		return nil, err
	}
	if !emailRegex.MatchString(email) {
		log.Warn().Str("email", email).Msg("Invalid email")
		return nil, errors.New("Invalid email")
	}
	return &User{Email: email, Password: password}, nil
}
