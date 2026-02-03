package testutil

import (
	"github.com/devbydaniel/announcable/internal/logger"
	"github.com/rs/zerolog"
)

var log = logger.Get()

// DisableLogging disables logging for tests
func DisableLogging() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// EnableLogging enables logging for tests at specified level
func EnableLogging(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// SetupTest performs common test setup and returns a cleanup function
func SetupTest() func() {
	// Disable logging by default for cleaner test output
	DisableLogging()

	// Return cleanup function
	return func() {
		// Reset logging to default
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
