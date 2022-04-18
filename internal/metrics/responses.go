package metrics

import "github.com/prometheus/client_golang/prometheus"

var ResponseCodes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_codes",
		Help: "Breakdown of aggregate HTTP response codes.",
	},
	[]string{"code"},
)

func init() {
	prometheus.Register(ResponseCodes)
}
