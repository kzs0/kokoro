package logs

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Logs struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogTimestampUnit string `env:"LOG_TIMESTAMP_UNIT" envDefault:"MILLISECONDS"`
	Pretty           bool   `env:"PRETTY_LOGS" envDefault:"false"`
	ServiceName      string `env:"SERVICE_NAME" envDefault:"_"`
	Environment      string `env:"ENVIRONMENT" envDefault:"dev"`
}

func parse(level string) (zerolog.Level, error) {
	return zerolog.ParseLevel(strings.ToLower(level))
}

func Init(config Logs) error {
	level, err := parse(config.LogLevel)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	switch strings.ToUpper(config.LogTimestampUnit) {
	case "SECONDS", "SECS", "SEC":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	case "MILLISECONDS", "MILLIS", "MILLI":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	case "MICROSECONDS", "MICROS", "MICRO":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	case "NANOSECONDS", "NANOS", "NANO":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	default:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	}

	if config.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Logger = log.With().Timestamp().Caller().
		Str("service", config.ServiceName).
		Str("env", config.Environment).
		Logger()

	return nil
}
