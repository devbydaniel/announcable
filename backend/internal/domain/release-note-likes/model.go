package releasenotelikes

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/google/uuid"
)

// ReleaseNoteLike represents a like on a release note by a client.
type ReleaseNoteLike struct {
	database.BaseModel `gorm:"embedded"`
	ReleaseNoteID      uuid.UUID
	ReleaseNote        releasenotes.ReleaseNote
	OrganisationID     uuid.UUID
	Organisation       organisation.Organisation
	ClientID           string `gorm:"type:text;not null"`
}
