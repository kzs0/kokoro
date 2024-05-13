package kokoro

import (
	"github.com/kenzo-spaulding/kokoro/telemetry/logs"
	"github.com/kenzo-spaulding/kokoro/telemetry/metrics"
	"github.com/kenzo-spaulding/kokoro/telemetry/traces"
)

type Config struct {
	logs.Logs
	metrics.Metrics
	traces.Traces
}
