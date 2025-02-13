package handler

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/devbydaniel/release-notes-go/internal/objstore"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog"
)

type Handler struct {
	DB       *database.DB
	ObjStore *objstore.ObjStore
	log      *zerolog.Logger
	decoder  *schema.Decoder
}

func NewHandler(db *database.DB, objStore *objstore.ObjStore) *Handler {
	log := logger.Get()
	return &Handler{DB: db, ObjStore: objStore, log: &log, decoder: schema.NewDecoder()}
}
