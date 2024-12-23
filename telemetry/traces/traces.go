package traces

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	api "go.opentelemetry.io/otel/sdk/trace"
)

// TODO endpoint for pushing traces and whether to use stdouttrace
type Traces struct {
	Enabled bool   `env:"TRACES_ENABLED" envDefault:"true"`
	Style   string `env:"TRACES_EXPORTER" envDefault:"CONSOLE"`
}

func Init(ctx context.Context, config Traces) error {
	if !config.Enabled {
		return nil
	}

	var exporter api.SpanExporter
	var err error

	switch strings.ToUpper(config.Style) {
	case "CONSOLE":
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		// Default No Export
	}

	if err != nil {
		return fmt.Errorf("failed to load trace exporter: %w", err)
	}

	bsp := api.NewBatchSpanProcessor(exporter)
	provider := api.NewTracerProvider(
		api.WithSampler(api.AlwaysSample()),
		api.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(provider)

	go func() {
		select {
		case <-ctx.Done():
			err = provider.Shutdown(ctx)
			if err != nil {
				slog.Error("faield to shutdown trace provider",
					slog.String("error", err.Error()))
			}
		}
	}()

	return nil
}
