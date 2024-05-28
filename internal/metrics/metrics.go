package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Devices       prometheus.Gauge
	Info          *prometheus.GaugeVec
	TotalRequests prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "urlshortener",
			Name:      "connected_devices",
			Help:      "The number of connected devices",
		}),
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "urlshortener",
			Name:      "info",
			Help:      "url-shortener version",
		}, []string{"version"}),
		TotalRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "urlshortener",
			Name:      "total_requests",
			Help:      "The number of requests",
		}),
	}

	m.Devices.Set(1)

	m.Info.With(prometheus.Labels{"version": "1.0.0"}).Set(1)

	reg.MustRegister(m.Devices, m.Info, m.TotalRequests)

	return m
}
