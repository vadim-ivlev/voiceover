package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NoColor — флаг для отключения цветного лога
var NoColor bool = false

// инициализируем логгер
func InitializeLogger() {
	// цветной лог в консоль
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: NoColor})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Logger initialized")
}
