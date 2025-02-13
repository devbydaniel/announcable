package mw

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/rs/zerolog"
)

type Handler struct {
	DB  *database.DB
	log zerolog.Logger
}

var log = logger.Get()

func NewHandler(db *database.DB) *Handler {
	return &Handler{DB: db, log: log}
}
