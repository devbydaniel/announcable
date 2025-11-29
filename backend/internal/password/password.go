package password

import (
	"errors"
	"regexp"

	"github.com/devbydaniel/announcable/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

var log = logger.Get()

func HashPassword(password string) (string, error) {
	log.Trace().Msg("HashPassword")
	passwordBytes := []byte(password)
	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func DoPasswordsMatch(hashedPassword, test string) bool {
	log.Trace().Msg("DoPasswordsMatch")
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(test)); err != nil {
		return false
	}
	return true
}

func IsValidPassword(password string) error {
	log.Trace().Msg("IsValidPassword")
	// Check minimum length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	uppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !uppercase {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	lowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !lowercase {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one number
	number := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !number {
		return errors.New("password must contain at least one number")
	}

	return nil
}
