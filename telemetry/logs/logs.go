package logs

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Logs struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	Pretty      bool   `env:"PRETTY_LOGS" envDefault:"false"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"_"`
	Environment string `env:"ENVIRONMENT" envDefault:"dev"`
}

var (
	ErrInitFailed  = errors.New("failed to initialize logs")
	ErrBadLogLevel = errors.New("invalid log level")
)

// Determines the log level from a provided string
// The string is trimmed of whitespaced and converted to uppercase
func ParseLevel(level string) (slog.Level, error) {
	switch strings.TrimSpace(strings.ToUpper(level)) {
	case "TRACE":
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
	case "FATAL":
	case "PANIC":
		return slog.LevelError, nil
	default:
	}

	err := fmt.Errorf("%s is not a valid log level", level)
	return slog.LevelInfo, errors.Join(ErrBadLogLevel, err)
}

func Init(config Logs) error {
	level, err := ParseLevel(config.LogLevel)
	if err != nil {
		return errors.Join(ErrInitFailed, err)
	}

	opts := slog.HandlerOptions{AddSource: true}
	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, &opts)

	if config.Pretty {
		handler = slog.NewTextHandler(os.Stdout, &opts)
	}

	defaultAttrs := []slog.Attr{
		slog.String("environment", config.Environment),
		slog.String("service", config.ServiceName),
	}

	handler = handler.WithAttrs(defaultAttrs)
	logger := slog.New(handler)

	slog.SetLogLoggerLevel(level)
	slog.SetDefault(logger)

	return nil
}
