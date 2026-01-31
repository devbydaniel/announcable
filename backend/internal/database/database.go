package database

import (
	"fmt"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.Get()

// DB wraps a GORM database client.
type DB struct {
	Client            *gorm.DB
	ErrRecordNotFound error
}

// Transaction wraps a GORM database transaction.
type Transaction struct {
	Tx *gorm.DB
}

var conf = config.New()

// Commit commits the transaction.
func (tx *Transaction) Commit() {
	if err := tx.Tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Error committing transaction")
	}
}

// Rollback rolls back the transaction.
func (tx *Transaction) Rollback() {
	if err := tx.Tx.Rollback().Error; err != nil {
		log.Error().Err(err).Msg("Error rolling back transaction")
	}
}

// StartTransaction begins a new database transaction.
func (db *DB) StartTransaction() *Transaction {
	tx := db.Client.Begin()
	return &Transaction{Tx: tx}
}

// Connect establishes a connection to the PostgreSQL database.
func Connect() (*DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Berlin",
		conf.Postgres.Host,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Name,
		conf.Postgres.Port,
	)
	log.Debug().Str("dsn", dsn).Msg("Connecting to database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return nil, err
	}
	return &DB{Client: db, ErrRecordNotFound: gorm.ErrRecordNotFound}, nil
}

// Close closes the database connection.
func Close(db *DB) {
	postgres, err := db.Client.DB()
	if err != nil {
		panic("Error while closing DB: " + err.Error())
	}
	if err := postgres.Close(); err != nil {
		panic("Error while closing DB: " + err.Error())
	}
}
