package mw

import (
	"github.com/devbydaniel/announcable/internal/database"
	"github.com/devbydaniel/announcable/internal/logger"
	"github.com/rs/zerolog"
)

// Handler holds dependencies for middleware functions.
type Handler struct {
	DB  *database.DB
	log zerolog.Logger
}

var log = logger.Get()

// NewHandler creates a new middleware Handler with the given database connection.
func NewHandler(db *database.DB) *Handler {
	return &Handler{DB: db, log: log}
}
