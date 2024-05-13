package traces

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	api "go.opentelemetry.io/otel/sdk/trace"
)

// TODO endpoint for pushing traces and whether to use stdouttrace
type Traces struct {
}

func Init(ctx context.Context, config Traces) error {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}

	bsp := api.NewBatchSpanProcessor(exp)
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
				log.Error().Err(err).Msg("failed to shutdown trace provider")
			}
		}
	}()

	return nil
}
