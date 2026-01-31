package session

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/google/uuid"
)

// Session represents an authenticated user session with an expiration time.
type Session struct {
	database.BaseModel `gorm:"embedded"`
	UserID             uuid.UUID
	ExpiresAt          int64 // UnixMilli
	ExternalID         string
}

// AuthCookieName is the name of the HTTP cookie used to store the session token.
var AuthCookieName = "announcable-session"

// New creates a new Session for the given user with the specified expiration and external session ID.
func New(userID uuid.UUID, expiresAt int64, sessionID string) Session {
	log.Trace().Str("userID", userID.String()).Int64("expiresAt", expiresAt).Str("sessionID", sessionID).Msg("New")
	return Session{UserID: userID, ExpiresAt: expiresAt, ExternalID: sessionID}
}
