package shared

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/devbydaniel/release-notes-go/internal/objstore"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog"
)

// Dependencies holds shared dependencies used across all handlers
type Dependencies struct {
	DB       *database.DB
	ObjStore *objstore.ObjStore
	Log      *zerolog.Logger
	Decoder  *schema.Decoder
}

// New creates a new Dependencies container with initialized dependencies
func New(db *database.DB, objStore *objstore.ObjStore) *Dependencies {
	log := logger.Get()
	return &Dependencies{
		DB:       db,
		ObjStore: objStore,
		Log:      &log,
		Decoder:  schema.NewDecoder(),
	}
}
