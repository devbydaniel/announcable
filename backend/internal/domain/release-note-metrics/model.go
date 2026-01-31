package releasenotemetrics

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/announcable/internal/domain/release-notes"
	"github.com/google/uuid"
)

// MetricType represents the type of release note metric being tracked.
type MetricType string

const (
	// MetricTypeView represents a release note view event.
	MetricTypeView MetricType = "view"
	// MetricTypeCtaClick represents a call-to-action click event.
	MetricTypeCtaClick MetricType = "cta_click"
)

// ReleaseNoteMetric represents a tracked metric event for a release note.
type ReleaseNoteMetric struct {
	database.BaseModel `gorm:"embedded"`
	ReleaseNoteID      uuid.UUID
	ReleaseNote        releasenotes.ReleaseNote
	OrganisationID     uuid.UUID
	Organisation       organisation.Organisation
	ClientID           string     `gorm:"type:text;not null"`
	MetricType         MetricType `gorm:"type:release_note_metric_type;not null"`
}
