package logger

import (
	"io"
	"os"

	adapter "github.com/axiomhq/axiom-go/adapters/zerolog"
	"github.com/devbydaniel/announcable/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	axiomWriter *adapter.Writer
	logger      zerolog.Logger
)

func init() {
	// Set global log level to trace
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	cfg := config.New()
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	var writers []io.Writer
	writers = append(writers, consoleWriter)

	// Only initialize Axiom if credentials are provided
	if cfg.Axiom.Token != "" && cfg.Axiom.Dataset != "" {
		var err error
		axiomWriter, err = adapter.New()
		if err != nil {
			// Log warning but don't crash - fall back to console only
			log.Warn().Err(err).Msg("Failed to create Axiom writer, using console logging only")
		} else {
			writers = append(writers, axiomWriter)
		}
	}

	logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).
		With().
		Timestamp().
		Caller().
		Str("env", cfg.Env).
		Logger()

	// Replace the global logger instance
	log.Logger = logger
}

// Get returns the configured zerolog.Logger instance
func Get() zerolog.Logger {
	return logger
}

// Cleanup properly closes the Axiom writer, ensuring all logs are flushed
// This should be called when your application is shutting down
func Cleanup() {
	if axiomWriter != nil {
		axiomWriter.Close()
	}
}
