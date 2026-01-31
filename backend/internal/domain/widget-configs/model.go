package widgetconfigs

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/domain/organisation"
	"github.com/google/uuid"
)

// WidgetConfig holds the appearance and behavior settings for an organisation's widget.
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
	ReleasePageBaseURL      *string    `gorm:"type:varchar(255)"`
	EnableLikes             bool       `gorm:"type:boolean;default:true"`
	LikeButtonText          string     `gorm:"type:varchar(255);default:'Like'"`
	UnlikeButtonText        string     `gorm:"type:varchar(255);default:'Unlike'"`
}

// WidgetType defines the display mode of the widget (modal, popover, or sidebar).
type WidgetType string

func (wt WidgetType) String() string {
	return string(wt)
}

// WidgetType constants define the available widget display modes.
const (
	WidgetTypeModal   WidgetType = "modal"
	WidgetTypePopover WidgetType = "popover"
	WidgetTypeSidebar WidgetType = "sidebar"
)
