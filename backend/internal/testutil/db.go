package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/devbydaniel/announcable/internal/database"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDBContainer holds the test database container and connection
type TestDBContainer struct {
	Container testcontainers.Container
	DB        *database.DB
	DSN       string
}

// SetupTestDB creates a PostgreSQL container and returns a test database connection
func SetupTestDB(t *testing.T) *TestDBContainer {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(60 * time.Second).
			WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable TimeZone=UTC",
		host, port.Port())

	// Connect to database
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db := &database.DB{
		Client:            gormDB,
		ErrRecordNotFound: gorm.ErrRecordNotFound,
	}

	return &TestDBContainer{
		Container: container,
		DB:        db,
		DSN:       dsn,
	}
}

// Cleanup terminates the test database container
func (tdb *TestDBContainer) Cleanup(t *testing.T) {
	ctx := context.Background()
	if err := tdb.Container.Terminate(ctx); err != nil {
		t.Errorf("Failed to terminate container: %v", err)
	}
}

// TruncateTable removes all data from a table
func (tdb *TestDBContainer) TruncateTable(tableName string) error {
	return tdb.DB.Client.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error
}

// TruncateAllTables removes all data from common tables
func (tdb *TestDBContainer) TruncateAllTables() error {
	tables := []string{
		"users",
		"organisations",
		"sessions",
		"release_notes",
		"release_note_likes",
		"release_note_metrics",
		"release_page_configs",
		"widget_configs",
	}

	for _, table := range tables {
		if err := tdb.TruncateTable(table); err != nil {
			log.Error().Err(err).Str("table", table).Msg("Failed to truncate table")
			return err
		}
	}

	return nil
}

// CreateTestID generates a new UUID for testing
func CreateTestID() uuid.UUID {
	return uuid.New()
}
