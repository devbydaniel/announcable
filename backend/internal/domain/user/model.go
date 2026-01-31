package user

import (
	"errors"
	"regexp"

	"github.com/devbydaniel/announcable/internal/database"
)

// User represents a registered user account with email and password credentials.
type User struct {
	database.BaseModel `gorm:"embedded"`
	Email              string `gorm:"unique"`
	Password           string
	EmailVerified      bool `gorm:"default:false"`
}

// New creates a new User after validating the email format.
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
