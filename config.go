package kokoro

import (
	"github.com/kzs0/kokoro/telemetry/logs"
	"github.com/kzs0/kokoro/telemetry/metrics"
	"github.com/kzs0/kokoro/telemetry/traces"
)

type Config struct {
	logs.Logs
	metrics.Metrics
	traces.Traces
}
