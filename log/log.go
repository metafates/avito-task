package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger = newLogger()

func newLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	consoleWriter := zerolog.NewConsoleWriter()

	return log.Output(consoleWriter)
}
