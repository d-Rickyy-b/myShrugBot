package prometheus

import (
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

var (
	ShrugUsage = metrics.NewCounter("myshrugbot_usage_total")
)

func StartPrometheusExporter(addr string) {
	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	http.ListenAndServe(addr, nil)
}
