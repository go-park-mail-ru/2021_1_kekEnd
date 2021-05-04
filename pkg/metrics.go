package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Hits    *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}

func NewMetrics(router *gin.Engine) *Metrics {
	var metrics Metrics

	metrics.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits",
		}, []string{"status", "path", "method"})

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "timings",
		}, []string{"status", "path", "method"})

	prometheus.MustRegister(metrics.Hits, metrics.Timings)

	router.Handle("/metrics", promhttp.Handler())

	return &metrics
}
