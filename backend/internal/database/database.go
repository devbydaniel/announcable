package database

import (
	"fmt"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.Get()

type DB struct {
	Client            *gorm.DB
	ErrRecordNotFound error
}

type Transaction struct {
	Tx *gorm.DB
}

var conf = config.New()

func (tx *Transaction) Commit() {
	if err := tx.Tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Error committing transaction")
	}
}

func (tx *Transaction) Rollback() {
	if err := tx.Tx.Rollback().Error; err != nil {
		log.Error().Err(err).Msg("Error rolling back transaction")
	}
}

func (db *DB) StartTransaction() *Transaction {
	tx := db.Client.Begin()
	return &Transaction{Tx: tx}
}

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

func Close(db *DB) {
	postgres, err := db.Client.DB()
	if err != nil {
		panic("Error while closing DB: " + err.Error())
	}
	if err := postgres.Close(); err != nil {
		panic("Error while closing DB: " + err.Error())
	}
}

func createDbIfNotExist(db *DB) error {
	// Check if the database exists
	var exists bool
	err := db.Client.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", conf.Postgres.Name).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create the database if it doesn't exist
	if !exists {
		createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %s", conf.Postgres.Name)
		if err := db.Client.Exec(createDatabaseCommand).Error; err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Printf("Database '%s' created successfully\n", conf.Postgres.Name)
	}
	if exists {
		fmt.Printf("Database '%v' already exists, skipping...", conf.Postgres.Name)
	}

	return nil
}
