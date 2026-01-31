package releasepageconfig

import (
	"io"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

// ReleasePageConfig holds the configuration for an organisation's public release page.
type ReleasePageConfig struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation       organisation.Organisation
	ImageURL           string `gorm:"-"` // only returned in GetReleaseNoteByID
	ImagePath          string `gorm:"type:varchar(255)"`
	Title              string `gorm:"type:varchar(255)"`
	Description        string `gorm:"type:text"`
	BgColor            string `gorm:"type:varchar(255)"`
	TextColor          string `gorm:"type:varchar(255)"`
	TextColorMuted     string `gorm:"type:varchar(255)"`
	BrandPosition      string `gorm:"type:varchar(255)"`
	BackLinkLabel      string `gorm:"type:varchar(255)"`
	BackLinkURL        string `gorm:"type:varchar(255)"`
	Slug               string `gorm:"type:varchar(255)"`
	DisableReleasePage bool   `gorm:"type:boolean;default:false"`
}

// BrandPosition defines the position of the brand image on the release page.
type BrandPosition string

func (bp BrandPosition) String() string {
	return string(bp)
}

// BrandPosition constants define the available brand image positions.
const (
	BrandPositionTop  BrandPosition = "top"
	BrandPositionLeft BrandPosition = "left"
)

// ImageInput holds image upload data for creating or updating a release page brand image.
type ImageInput struct {
	ShouldDeleteImage bool
	ImgData           io.Reader
	Format            string
}
