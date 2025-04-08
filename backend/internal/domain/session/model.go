package session

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/google/uuid"
)

type Session struct {
	database.BaseModel `gorm:"embedded"`
	UserID             uuid.UUID
	ExpiresAt          int64 // UnixMilli
	ExternalID         string
}

var AuthCookieName = "announcable-session"

func New(userId uuid.UUID, expiresAt int64, sessionId string) Session {
	log.Trace().Str("userId", userId.String()).Int64("expiresAt", expiresAt).Str("sessionId", sessionId).Msg("New")
	return Session{UserID: userId, ExpiresAt: expiresAt, ExternalID: sessionId}
}
