package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Metrics map[string]prometheus.Counter

func StartPrometheusHttpServer(port string, network_name string) {
	Metrics = make(map[string]prometheus.Counter)

	// Create gauges
	Metrics["kafka_messages_consumed"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "kafka_messages_consumed",
		Help:        "number of messages read by the consumer",
		ConstLabels: prometheus.Labels{"network_name": network_name},
	})

	Metrics["websockets_connected"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "websockets_connected",
		Help:        "number of websockets connected",
		ConstLabels: prometheus.Labels{"network_name": network_name},
	})

	Metrics["websockets_bytes_written"] = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "websockets_bytes_written",
		Help:        "number of bytes sent over through websockets",
		ConstLabels: prometheus.Labels{"network_name": network_name},
	})

	// Start server
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+port, nil)
}
