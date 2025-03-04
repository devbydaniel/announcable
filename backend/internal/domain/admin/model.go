package admin

import (
	"github.com/google/uuid"
)

// AdminUser represents the admin user with special privileges
type AdminUser struct {
	UserID uuid.UUID
}

// IsAdmin checks if the provided user ID matches the admin user ID
func IsAdmin(userID uuid.UUID, adminUserID string) bool {
	if adminUserID == "" {
		return false
	}

	adminID, err := uuid.Parse(adminUserID)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing admin user ID")
		return false
	}

	return userID == adminID
}
