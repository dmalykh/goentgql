package middleware

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

func Metrics() graphql.OperationMiddleware {
	var totalRequestMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_total",
		Help: "The total number of request",
	})
	var operationDurationMetric = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_server_request_duration_seconds",
		Help:    "Histogram of response time for handler in seconds",
		Buckets: []float64{.000001, .00001, .0001, .001, .01, .025, .05, .1, .25, .5, 1},
	})

	return func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
		totalRequestMetric.Inc()
		var start = time.Now()
		var resp = next(ctx)
		var t = time.Since(start).Seconds()
		operationDurationMetric.Observe(t)
		return resp
	}
}
