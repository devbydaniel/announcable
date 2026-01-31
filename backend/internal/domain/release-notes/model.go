package releasenotes

import (
	"io"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

// ReleaseNote represents a release note entry belonging to an organisation.
type ReleaseNote struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation       organisation.Organisation
	Title              string             `gorm:"type:varchar(255)"`
	ImageURL           string             `gorm:"-"` // only returned in GetReleaseNoteByID
	DescriptionShort   string             `gorm:"type:text"`
	DescriptionLong    string             `gorm:"type:text"`
	ReleaseDate        *string            `gorm:"type:date;default:null"`
	ImagePath          string             `gorm:"type:varchar(255)"`
	MediaLink          string             `gorm:"type:varchar(1024)"`
	IsPublished        bool               `gorm:"type:bool;default:false"`
	CtaLabelOverride   string             `gorm:"type:varchar(255)"`
	CtaURLOverride     string             `gorm:"type:varchar(255)"`
	HideCta            bool               `gorm:"type:bool;default:false"`
	AttentionMechanism AttentionMechanism `gorm:"type:varchar(255)"`
	CreatedBy          uuid.UUID          `gorm:"type:uuid"`
	LastUpdatedBy      uuid.UUID          `gorm:"type:uuid"`
	HideOnWidget       bool               `gorm:"type:bool;default:false"`
	HideOnReleasePage  bool               `gorm:"type:bool;default:false"`
}

// PaginatedReleaseNotes holds a page of release notes along with pagination metadata.
type PaginatedReleaseNotes struct {
	Items      []*ReleaseNote
	TotalCount int64
	TotalPages int
	Page       int
	PageSize   int
}

// ReleaseNoteStatus holds the update timestamp and attention mechanism for a published release note.
type ReleaseNoteStatus struct {
	UpdatedAt          string
	AttentionMechanism string
}

// AttentionMechanism defines how a release note attracts user attention in the widget.
type AttentionMechanism string

func (am AttentionMechanism) String() string {
	return string(am)
}

// AttentionMechanism constants define the available attention mechanisms.
const (
	AttentionMechanismIndicator   AttentionMechanism = "show_indicator"
	AttentionMechanismInstantOpen AttentionMechanism = "instant_open"
)

// ImageInput holds image upload data for creating or updating a release note image.
type ImageInput struct {
	ShouldDeleteImage bool
	ImgData           io.Reader
	Format            string
}
