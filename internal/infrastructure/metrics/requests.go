package metrics

import "github.com/prometheus/client_golang/prometheus"

var TotalHTTPRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rpc_http_requests_total",
		Help: "Number of rpc http requests.",
	},
	[]string{"method"},
)

var TotalWSRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rpc_ws_requests_total",
		Help: "Number of rpc ws requests.",
	},
	[]string{"method"},
)

func init() {
	prometheus.Register(TotalHTTPRequests)
	prometheus.Register(TotalWSRequests)
}
