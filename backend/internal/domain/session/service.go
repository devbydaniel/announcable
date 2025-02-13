package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const sessionExpiresIn = 30 * 24 * time.Hour

type service struct {
	repository repository
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repository: r}
}

func (s *service) CreateToken() string {
	log.Trace().Msg("CreateSessionToken")
	bytes := make([]byte, 15)
	rand.Read(bytes)
	token := base32.StdEncoding.EncodeToString(bytes)
	log.Debug().Str("token", token).Msg("")
	return token
}

func (s *service) Create(token string, userId uuid.UUID) error {
	log.Trace().Str("token", token).Str("userId", userId.String()).Msg("CreateSession")
	sessionId := getIdFromToken(token)
	expiresAt := calcNextExpiry()
	session := Session{ExternalID: sessionId, ExpiresAt: expiresAt, UserID: userId}
	return s.repository.Save(&session)
}

func (s *service) CreateCustomDuration(token string, userId uuid.UUID, duration time.Duration) error {
	log.Trace().Str("token", token).Str("userId", userId.String()).Msg("CreateCustomDuration")
	sessionId := getIdFromToken(token)
	expiresAt := time.Now().Add(duration).UnixMilli()
	session := Session{ExternalID: sessionId, ExpiresAt: expiresAt, UserID: userId}
	return s.repository.Save(&session)
}

func (s *service) ValidateSession(token string) (*Session, error) {
	log.Trace().Str("token", token).Msg("ValidateSession")
	sessionId := getIdFromToken(token)
	session, err := s.repository.FindByExternalId(sessionId)
	if err != nil {
		return nil, err
	}
	if sessionIsExpired(session) {
		log.Debug().Msg("session expired")
		if err := s.repository.Delete(session.ID); err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		return nil, s.repository.db.ErrRecordNotFound
	}
	session.ExpiresAt = calcNextExpiry()
	if err := s.repository.Save(session); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return session, nil
}

func (s *service) InvalidateUserSessions(userId uuid.UUID) error {
	log.Trace().Str("userId", userId.String()).Msg("InvalidateUserSessions")
	return s.repository.DeleteByUserId(userId)
}

func (s *service) Delete(id uuid.UUID) error {
	log.Trace().Str("token", id.String()).Msg("DeleteSession")
	return s.repository.Delete(id)
}

func getIdFromToken(token string) string {
	byteId := sha256.Sum256([]byte(token))
	return hex.EncodeToString(byteId[:])
}

func calcNextExpiry() int64 {
	return time.Now().Add(sessionExpiresIn).UnixMilli()
}

func sessionIsExpired(s *Session) bool {
	return time.Now().After(time.UnixMilli(s.ExpiresAt))
}
