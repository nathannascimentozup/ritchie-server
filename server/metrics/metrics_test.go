package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestMetric(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		in   args
		out  string
	}{
		{
			name: "tree",
			in:   args{path: "/tree"},
			out:  `Desc{fqName: "http_request_tree", help: "The total number service calls to path /tree", constLabels: {}, variableLabels: [code]}`,
		},
		{
			name: "health",
			in:   args{path: "/health"},
			out:  `Desc{fqName: "http_request_health", help: "The total number service calls to path /health", constLabels: {}, variableLabels: [code]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Metric(tt.in.path).GetMetricWith(prometheus.Labels{"code": "200"}); fmt.Sprint(got.Desc()) != tt.out {
				t.Errorf("got = %v, want %v", got.Desc(), tt.out)
			}
		})
	}
}
