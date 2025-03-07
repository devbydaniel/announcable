package logger

import (
	"os"

	adapter "github.com/axiomhq/axiom-go/adapters/zerolog"
	"github.com/devbydaniel/release-notes-go/config"
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
	
	env := config.New().Env
	var err error
	axiomWriter, err = adapter.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Axiom writer")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	logger = zerolog.New(zerolog.MultiLevelWriter(axiomWriter, consoleWriter)).
		With().
		Timestamp().
		Caller().
		Str("env", env).
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
