package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// const (
// 	disconnected = 0.0
// 	connected    = 1.0

// 	// traceIDKey is used as the trace ID key value in the prometheus.Labels in a prometheus.Exemplar.
// 	//
// 	// Its value of `trace_id` complies with the OpenTelemetry specification for metrics' exemplars, as seen in:
// 	// https://opentelemetry.io/docs/specs/otel/metrics/data-model/#exemplars
// 	traceIDKey = "trace_id"
// )

type Metrics struct {
	httpSendCount *prometheus.CounterVec
	collectors    []prometheus.Collector
}

func NewMetrics() *Metrics {
	//
	// You need to add the metric also in the Registry() method
	//
	return &Metrics{
		httpSendCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_send_count",
			Help: "Number HTTP send",
		}, []string{"source"}),
	}
}

func (m *Metrics) IncHttpSend(source string) {
	m.httpSendCount.WithLabelValues(source).Inc()
}

func (m *Metrics) Registry() (*prometheus.Registry, error) {
	reg := prometheus.NewRegistry()

	for _, metric := range []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			ReportErrors: false,
		}),
		//
		// add custom metrics here to be exposed by endpoint
		//
		m.httpSendCount,
	} {
		err := reg.Register(metric)
		if err != nil {
			return nil, err
		}
	}

	for _, metric := range m.collectors {
		err := reg.Register(metric)
		if err != nil {
			return nil, err
		}
	}

	return reg, nil
}
