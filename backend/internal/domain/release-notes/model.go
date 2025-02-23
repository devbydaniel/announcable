package releasenotes

import (
	"io"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/google/uuid"
)

type ReleaseNote struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation       organisation.Organisation
	Title              string             `gorm:"type:varchar(255)"`
	ImageUrl           string             `gorm:"-"` // only returned in GetReleaseNoteByID
	DescriptionShort   string             `gorm:"type:text"`
	DescriptionLong    string             `gorm:"type:text"`
	ReleaseDate        *string            `gorm:"type:date;default:null"`
	ImagePath          string             `gorm:"type:varchar(255)"`
	IsPublished        bool               `gorm:"type:bool;default:false"`
	CtaLabelOverride   string             `gorm:"type:varchar(255)"`
	CtaUrlOverride     string             `gorm:"type:varchar(255)"`
	HideCta            bool               `gorm:"type:bool;default:false"`
	AttentionMechanism AttentionMechanism `gorm:"type:varchar(255)"`
	CreatedBy          uuid.UUID          `gorm:"type:uuid"`
	LastUpdatedBy      uuid.UUID          `gorm:"type:uuid"`
	HideOnWidget       bool               `gorm:"type:bool;default:false"`
	HideOnReleasePage  bool               `gorm:"type:bool;default:false"`
}

type PaginatedReleaseNotes struct {
	Items      []*ReleaseNote
	TotalCount int64
	TotalPages int
	Page       int
	PageSize   int
}

type ReleaseNoteStatus struct {
	UpdatedAt          string
	AttentionMechanism string
}

type AttentionMechanism string

func (am AttentionMechanism) String() string {
	return string(am)
}

const (
	AttentionMechanismIndicator   AttentionMechanism = "show_indicator"
	AttentionMechanismInstantOpen AttentionMechanism = "instant_open"
)

type ImageInput struct {
	ShoudDeleteImage bool
	ImgData          io.Reader
	Format           string
}
