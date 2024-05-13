package logs

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
)

type Logs struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogTimestampUnit string `env:"LOG_TIMESTAMP_UNIT" envDefault:"milliseconds"`
	Pretty           bool   `env:"PRETTY_LOGS" envDefault:"false"`
	ServiceName      string `env:"SERVICE_NAME" envDefault:"-"`
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

	switch strings.ToLower(config.LogTimestampUnit) {
	case "seconds", "secs", "sec":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	case "milliseconds", "millis", "milli":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	case "microseconds", "micros", "micro":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	case "nanoseconds", "nanos", "nano":
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
