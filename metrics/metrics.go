package metrics

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/gramework/gramework"
	"github.com/gramework/runtimer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Middleware handles metrics data
type Middleware struct {
	httpReqCounter *prometheus.CounterVec
	reqDuration    *prometheus.HistogramVec
}

const (
	typeHTTPS = "https"
	typeHTTP  = "http"
)

var metricsPath = []byte("/metrics")

// Register the middlewares
func Register(app *gramework.App, serviceName ...string) error {
	var m Middleware
	name := os.Args[0]
	if len(serviceName) > 0 {
		name = serviceName[0]
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	m.httpReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gramework_http_requests_total",
			Help: "Total count of HTTP requests processed, partitioned by code, method, path and type (HTTP/HTTPS)",
			ConstLabels: prometheus.Labels{
				"service": name,
				"node":    hostname,
			},
		},
		[]string{"code", "method", "path", "type"},
	)
	if err = prometheus.Register(m.httpReqCounter); err != nil {
		return err
	}

	m.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "gramework_http_requests_duration_seconds",
			Help: "Request processing duration, partitioned by code, method, path and type (HTTP/HTTPS)",
			ConstLabels: prometheus.Labels{
				"service": name,
				"node":    hostname,
			},
		},
		[]string{"code", "method", "path", "type"},
	)

	if err = prometheus.Register(m.reqDuration); err != nil {
		return err
	}

	app.GET("/metrics", gramework.NewGrameHandler(promhttp.Handler()))

	if err = app.UsePre(m.startReq); err != nil {
		return err
	}

	return app.UseAfterRequest(m.endReq)
}

func (m *Middleware) startReq(ctx *gramework.Context) {
	if bytes.Equal(ctx.Path(), metricsPath) {
		return
	}
	ctx.SetUserValue("gramework.metrics.startTime", time.Now())
}

func (m *Middleware) endReq(ctx *gramework.Context) {
	if bytes.Equal(ctx.Path(), metricsPath) {
		return
	}

	opts := []string{
		fmt.Sprintf("%d", ctx.Response.StatusCode()),
		gramework.BytesToString(ctx.Method()),
		gramework.BytesToString(ctx.Path()),
		"",
	}

	if ctx.IsTLS() {
		opts[3] = typeHTTPS
	} else {
		opts[3] = typeHTTP
	}

	m.httpReqCounter.WithLabelValues(opts...).Add(1)

	startTime := float64(
		time.Since(
			*(*time.Time)(runtimer.GetEfaceDataPtr(ctx.UserValue("gramework.metrics.startTime"))),
		).Nanoseconds(),
	) / float64(time.Millisecond)

	m.reqDuration.WithLabelValues(opts...).Observe(startTime)
}
