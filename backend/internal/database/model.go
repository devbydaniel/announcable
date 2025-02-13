package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model `gorm:"embedded"`
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
}
