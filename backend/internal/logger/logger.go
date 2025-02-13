package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Get() zerolog.Logger {
	return log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
