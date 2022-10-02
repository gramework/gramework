package metrics

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"github.com/gramework/gramework"
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

	millisecond = float64(time.Millisecond)

	uvKey = "gramework.metrics.startTime"
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

	app.GET(string(metricsPath), gramework.NewGrameHandler(promhttp.Handler()))

	if err = app.UsePre(m.startReq); err != nil {
		return err
	}

	return app.UseAfterRequest(m.endReq)
}

func (m *Middleware) startReq(ctx *gramework.Context) {
	if bytes.Equal(ctx.Path(), metricsPath) {
		return
	}
	ctx.SetUserValue(uvKey, gramework.Nanotime())
}

func (m *Middleware) endReq(ctx *gramework.Context) {
	if bytes.Equal(ctx.Path(), metricsPath) {
		return
	}

	opts := []string{
		strconv.FormatInt(int64(ctx.Response.StatusCode()), 10),
		string(ctx.Method()),
		string(ctx.Path()),
		typeHTTP,
	}

	if ctx.IsTLS() {
		opts[3] = typeHTTPS
	}

	m.httpReqCounter.WithLabelValues(opts...).Add(1)

	startTime, _ := ctx.UserValue(uvKey).(int64)
	duration := float64(gramework.Nanotime()-startTime) / millisecond

	m.reqDuration.WithLabelValues(opts...).Observe(duration)
}
