package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
)

var (
	metrics           = make(map[string]*prometheus.CounterVec)
	LatencyOpsRequest = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_latency",
		Help:    "The http service latency",
		Buckets: []float64{0, 10, 50, 100, 200, 300, 500, 1000, 2000, 5000},
	},
		[]string{
			"path",
		})
)

func Metric(path string) *prometheus.CounterVec {
	metric := metrics[path]
	if metric == nil {
		name := strings.ReplaceAll(path, "/", "_")
		name = strings.ReplaceAll(name, "-", "_")
		metric = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request" + name,
			Help: "The total number service calls to path " + path,
		},
			[]string{
				"code",
			})
		metrics[path] = metric
	}
	return metric
}
