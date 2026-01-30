package releasepageconfig

import (
	"io"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

type ReleasePageConfig struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation       organisation.Organisation
	ImageUrl           string `gorm:"-"` // only returned in GetReleaseNoteByID
	ImagePath          string `gorm:"type:varchar(255)"`
	Title              string `gorm:"type:varchar(255)"`
	Description        string `gorm:"type:text"`
	BgColor            string `gorm:"type:varchar(255)"`
	TextColor          string `gorm:"type:varchar(255)"`
	TextColorMuted     string `gorm:"type:varchar(255)"`
	BrandPosition      string `gorm:"type:varchar(255)"`
	BackLinkLabel      string `gorm:"type:varchar(255)"`
	BackLinkUrl        string `gorm:"type:varchar(255)"`
	Slug               string `gorm:"type:varchar(255)"`
	DisableReleasePage bool   `gorm:"type:boolean;default:false"`
}

type BrandPosition string

func (bp BrandPosition) String() string {
	return string(bp)
}

const (
	BrandPositionTop  BrandPosition = "top"
	BrandPositionLeft BrandPosition = "left"
)

type ImageInput struct {
	ShouldDeleteImage bool
	ImgData           io.Reader
	Format            string
}
