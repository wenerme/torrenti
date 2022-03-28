package serves

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func RegisterMetrics() {
	RegisterEndpoints(&HTTPEndpoint{
		Path:    "/metrics",
		Handler: promhttp.Handler(),
		EndpointDesc: EndpointDesc{
			Name:        "Metrics",
			Description: "Metrics exposed by the server",
			Selector:    "debug",
		},
	})
}

type MetricsMiddleware struct {
	// Prefix is the prefix that will be set on the metrics, by default it will be empty.
	Prefix string
	// DurationBuckets are the buckets used by Prometheus for the HTTP request duration metrics,
	// by default uses Prometheus default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// SizeBuckets are the buckets used by Prometheus for the HTTP response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// Registry is the registry that will be used by the recorder to store the metrics,
	// if the default registry is not used then it will use the default one.
	Registry prometheus.Registerer
	// HandlerLabel is the name that will be set to the handler ID label, by default is `handler`.
	HandlerLabel string
	// ServiceLabel is the name that will be set to the service label, by default is `service`.
	ServiceLabel string
	// ConstLabels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their ConstLabels.
	ConstLabels prometheus.Labels
}

func (mm *MetricsMiddleware) defaults() {
	if len(mm.DurationBuckets) == 0 {
		mm.DurationBuckets = prometheus.DefBuckets
	}

	if len(mm.SizeBuckets) == 0 {
		mm.SizeBuckets = prometheus.ExponentialBuckets(100, 10, 8)
	}

	if mm.Registry == nil {
		mm.Registry = prometheus.DefaultRegisterer
	}

	if mm.HandlerLabel == "" {
		mm.HandlerLabel = "handler"
	}

	if mm.ServiceLabel == "" {
		mm.ServiceLabel = "service"
	}
}

func (mm *MetricsMiddleware) Handle() func(next http.Handler) http.Handler {
	mm.defaults()

	httpRequestDurHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   mm.Prefix,
		Subsystem:   "http",
		Name:        "request_duration_seconds",
		Help:        "The latency of the HTTP requests.",
		Buckets:     mm.DurationBuckets,
		ConstLabels: mm.ConstLabels,
	}, []string{"code", "method"})

	httpRequestSizeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   mm.Prefix,
		Subsystem:   "http",
		Name:        "request_size_bytes",
		Help:        "The size of the HTTP requests.",
		Buckets:     mm.SizeBuckets,
		ConstLabels: mm.ConstLabels,
	}, []string{"code", "method"})

	httpResponseSizeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   mm.Prefix,
		Subsystem:   "http",
		Name:        "response_size_bytes",
		Help:        "The size of the HTTP responses.",
		Buckets:     mm.SizeBuckets,
		ConstLabels: mm.ConstLabels,
	}, []string{"code", "method"})

	httpRequestsInflight := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   mm.Prefix,
		Subsystem:   "http",
		Name:        "requests_inflight",
		Help:        "The number of inflight requests being handled at the same time.",
		ConstLabels: mm.ConstLabels,
	}, []string{}).With(prometheus.Labels{})

	httpRequestTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   mm.Prefix,
		Subsystem:   "http",
		Name:        "requests_total",
		Help:        "Total number of scrapes by HTTP status code.",
		ConstLabels: mm.ConstLabels,
	}, []string{"code"})

	for _, v := range []int{200, 400, 404, 500} {
		httpRequestTotal.WithLabelValues(strconv.Itoa(v))
	}

	log.Info().Msg("Registering metrics")
	mm.Registry.MustRegister(
		httpRequestTotal,
		httpRequestsInflight,
		httpRequestDurHistogram,
		httpResponseSizeHistogram,
		httpRequestSizeHistogram,
	)

	return func(next http.Handler) http.Handler {
		next = promhttp.InstrumentHandlerCounter(httpRequestTotal, next)
		next = promhttp.InstrumentHandlerInFlight(httpRequestsInflight, next)
		next = promhttp.InstrumentHandlerDuration(httpRequestDurHistogram, next)
		next = promhttp.InstrumentHandlerResponseSize(httpResponseSizeHistogram, next)
		next = promhttp.InstrumentHandlerRequestSize(httpRequestSizeHistogram, next)
		return next
	}
}
