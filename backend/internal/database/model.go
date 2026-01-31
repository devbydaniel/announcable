package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel provides common fields for all database models.
type BaseModel struct {
	gorm.Model `gorm:"embedded"`
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
}
