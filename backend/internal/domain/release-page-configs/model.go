package releasepageconfig

import (
	"io"

	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/google/uuid"
)

type ReleasePageConfig struct {
	database.BaseModel `gorm:"embedded"`
	OrganisationID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation       organisation.Organisation
	ImageUrl           string `gorm:"-"` // only returned in GetReleaseNoteByID
	Title              string `gorm:"type:varchar(255)"`
	Description        string `gorm:"type:text"`
	BgColor            string `gorm:"type:varchar(255)"`
	TextColor          string `gorm:"type:varchar(255)"`
	TextColorMuted     string `gorm:"type:varchar(255)"`
	BrandPosition      string `gorm:"type:varchar(255)"`
}

type BrandPosition string

func (bp BrandPosition) String() string {
	return string(bp)
}

const (
	BrandPositionTop   BrandPosition = "top"
	BrandPositionLeft  BrandPosition = "left"
	BrandPositionRight BrandPosition = "right"
)

type ImageInput struct {
	ShoudDeleteImage bool
	ImgData          io.Reader
	Format           string
}
