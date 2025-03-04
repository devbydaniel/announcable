package subscription

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/domain/organisation"
	"github.com/google/uuid"
)

type Subscription struct {
	database.BaseModel   `gorm:"embedded"`
	OrganisationID       uuid.UUID `gorm:"type:uuid"`
	Organisation         organisation.Organisation
	StripeSubscriptionID string `gorm:"type:text;not null"`
	IsActive             bool   `gorm:"type:boolean;not null"`
	IsFree               bool   `gorm:"type:boolean;default:false"`
}
