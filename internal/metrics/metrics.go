package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	devices prometheus.Gauge
	info    *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "urlshortener",
			Name:      "connected_devices",
			Help:      "The number of connected devices",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "urlshortener",
			Name:      "info",
			Help:      "url-shortener version",
		}, []string{"version"}),
	}

	m.devices.Set(1)
	m.info.With(prometheus.Labels{"version": "1.0.0"}).Set(1)

	reg.MustRegister(m.devices, m.info)

	return m
}
