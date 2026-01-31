package releasenotelikes

import (
	"github.com/google/uuid"
)

type service struct {
	repo *repository
}

// NewService creates a new release note likes service.
func NewService(r *repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

func (s *service) ToggleLike(releaseNoteID uuid.UUID, orgID uuid.UUID, clientID string) (bool, error) {
	log.Trace().Msg("ToggleLike")

	// Check if like already exists
	existingLike, err := s.repo.FindByReleaseNoteAndClientID(releaseNoteID, clientID)
	if err != nil {
		return false, err
	}

	// If like exists, delete it (unlike)
	if existingLike != nil {
		if err := s.repo.Delete(existingLike); err != nil {
			return true, err // Return true because it's still liked if delete failed
		}
		return false, nil // Successfully unliked
	}

	// If like doesn't exist, create it
	like := &ReleaseNoteLike{
		ReleaseNoteID:  releaseNoteID,
		OrganisationID: orgID,
		ClientID:       clientID,
	}
	if err := s.repo.Create(like); err != nil {
		return false, err
	}
	return true, nil // Successfully liked
}

func (s *service) GetLikeCount(releaseNoteID uuid.UUID) (int, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("GetLikeCount")
	likes, err := s.repo.FindByReleaseNoteID(releaseNoteID)
	if err != nil {
		return 0, err
	}
	return len(likes), nil
}

func (s *service) HasUserLiked(releaseNoteID uuid.UUID, clientID string) (bool, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Str("clientID", clientID).Msg("HasUserLiked")
	like, err := s.repo.FindByReleaseNoteAndClientID(releaseNoteID, clientID)
	if err != nil {
		return false, err
	}
	return like != nil, nil
}

func (s *service) GetLikesByReleaseNote(releaseNoteID uuid.UUID) ([]ReleaseNoteLike, error) {
	log.Trace().Str("releaseNoteID", releaseNoteID.String()).Msg("GetLikesByReleaseNote")
	return s.repo.FindByReleaseNoteID(releaseNoteID)
}

func (s *service) GetLikesByOrg(orgID uuid.UUID) ([]ReleaseNoteLike, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetLikesByOrg")
	return s.repo.FindByOrgID(orgID)
}
