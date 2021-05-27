package middleware

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
)

// AccessLog логгер
type AccessLog interface {
	AccessLogMiddleware(log *logger.Logger) gin.HandlerFunc
}

var (
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})

	timings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "method_timings",
		Help: "Per method timing",
	}, []string{"method"})
)

// AccessLogMiddleware мидлвара логгера
func AccessLogMiddleware(log *logger.Logger) gin.HandlerFunc {
	prometheus.MustRegister(hits, timings)

	return func(ctx *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		rid := fmt.Sprintf("%016x", rand.Int())[:5]

		log.StartReq(*ctx.Request, rid)
		start := time.Now()

		ctx.Next()

		hits.
			WithLabelValues(strconv.Itoa(ctx.Writer.Status()), ctx.Request.URL.String()).
			Inc()

		timings.
			WithLabelValues(ctx.Request.URL.String()).
			Observe(time.Since(start).Seconds())

		log.EndReq(*ctx.Request, start, rid)
	}
}
