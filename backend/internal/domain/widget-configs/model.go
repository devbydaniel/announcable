package widgetconfigs

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/google/uuid"
)

type WidgetConfig struct {
	database.BaseModel      `gorm:"embedded"`
	OrganisationID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Organisation            organisation.Organisation
	Title                   string     `gorm:"type:varchar(255)"`
	Description             string     `gorm:"type:text"`
	WidgetBorderRadius      int        `gorm:"type:int;default:12"`
	WidgetBorderColor       string     `gorm:"type:varchar(255)"`
	WidgetBorderWidth       int        `gorm:"type:int;default:1"`
	WidgetBgColor           string     `gorm:"type:varchar(255)"`
	WidgetTextColor         string     `gorm:"type:varchar(255)"`
	WidgetType              WidgetType `gorm:"type:varchar(255)"`
	ReleaseNoteBorderRadius int        `gorm:"type:int;default:8"`
	ReleaseNoteBorderColor  string     `gorm:"type:varchar(255)"`
	ReleaseNoteBorderWidth  int        `gorm:"type:int;default:0"`
	ReleaseNoteBgColor      string     `gorm:"type:varchar(255)"`
	ReleaseNoteTextColor    string     `gorm:"type:varchar(255)"`
	ReleaseNoteCtaText      string     `gorm:"type:varchar(255)"`
	ReleasePageBaseUrl      *string    `gorm:"type:varchar(255)"`
	EnableLikes             bool       `gorm:"type:boolean;default:true"`
	LikeButtonText          string     `gorm:"type:varchar(255);default:'Like'"`
	UnlikeButtonText        string     `gorm:"type:varchar(255);default:'Unlike'"`
}

type WidgetType string

func (wt WidgetType) String() string {
	return string(wt)
}

const (
	WidgetTypeModal   WidgetType = "modal"
	WidgetTypePopover WidgetType = "popover"
	WidgetTypeSidebar WidgetType = "sidebar"
)
