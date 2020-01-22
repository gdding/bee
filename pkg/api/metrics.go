package api

import (
	"net/http"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	// all metrics fields must be exported
	// to be able to return them by Metrics()
	// using reflection
	RequestCount     prometheus.Counter
	ResponseDuration prometheus.Histogram
	PingRequestCount prometheus.Counter
}

func newMetrics() (m metrics) {
	return metrics{
		RequestCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_request_count",
			Help: "Number of API requests.",
		}),
		ResponseDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "api_response_duration_seconds",
			Help:    "Histogram of API response durations.",
			Buckets: []float64{0.01, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		}),
		PingRequestCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_ping_request_count",
			Help: "Number HTTP API ping requests.",
		}),
	}
}

func (s *server) Metrics() (cs []prometheus.Collector) {
	v := reflect.Indirect(reflect.ValueOf(s.metrics))
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		if u, ok := v.Field(i).Interface().(prometheus.Collector); ok {
			cs = append(cs, u)
		}
	}
	return cs
}

func (s *server) pageviewMetricsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.metrics.RequestCount.Inc()
		h.ServeHTTP(w, r)
		s.metrics.ResponseDuration.Observe(time.Since(start).Seconds())
	})
}