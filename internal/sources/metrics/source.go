package metrics

import (
	"fmt"
	"runtime"
	"time"

	"go-project-template/internal/pkg/configs"

	"github.com/VictoriaMetrics/metrics"
)

func Init() {
	metrics.NewCounter(
		fmt.Sprintf(`version{version="%v",hostname="%v"}`,
			configs.Hostname, configs.Version),
	).Set(uint64(runtime.NumCPU()))
}

func PGRequest(query string, statTime time.Time) {
	metrics.GetOrCreateHistogram(fmt.Sprintf(`pg{query="%v"}`, query)).
		UpdateDuration(statTime)
}

func SchedulerChecks(scheduler, op string, counter int) {
	metrics.GetOrCreateCounter(fmt.Sprintf(`schedule{entity="%v",op="%v"}`,
		scheduler, op)).Add(counter)
}
