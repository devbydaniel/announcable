package releasenotemetrics

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	releasenotes "github.com/devbydaniel/release-notes-go/internal/domain/release-notes"
	"github.com/google/uuid"
)

type MetricType string

const (
	MetricTypeView     MetricType = "view"
	MetricTypeCtaClick MetricType = "cta_click"
)

type ReleaseNoteMetric struct {
	database.BaseModel `gorm:"embedded"`
	ReleaseNoteID      uuid.UUID
	ReleaseNote        releasenotes.ReleaseNote
	OrganisationID     uuid.UUID
	Organisation       organisation.Organisation
	ClientID           string     `gorm:"type:text;not null"`
	MetricType         MetricType `gorm:"type:release_note_metric_type;not null"`
}
