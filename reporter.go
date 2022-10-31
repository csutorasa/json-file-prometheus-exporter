package main

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsReporter struct {
	metricName string
	labelKeys  []string
	labels     map[string][]string
	counter    *prometheus.CounterVec
}

func NewReporter(metricName string) *MetricsReporter {
	return &MetricsReporter{
		metricName: metricName,
	}
}

func (r *MetricsReporter) Init(labels []string) {
	r.labels = map[string][]string{}
	r.labelKeys = []string{}
	for _, l := range labels {
		label := strings.ReplaceAll(l, ".", "_")
		r.labels[label] = strings.Split(l, ".")
		r.labelKeys = append(r.labelKeys, label)
	}
	logger.Printf("Using labels %v", r.labelKeys)
	r.counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: r.metricName,
	}, r.labelKeys)
}

func (r *MetricsReporter) Report(data map[string]any) {
	reportValues := []string{}
	for _, key := range r.labelKeys {
		reportValues = append(reportValues, getValue(data, r.labels[key]))
	}
	r.counter.WithLabelValues(reportValues...).Inc()
}

func getValue(data map[string]any, path []string) string {
	value, ok := data[path[0]]
	if !ok {
		return ""
	}
	if len(path) == 1 {
		return fmt.Sprintf("%v", value)
	}
	m, ok := value.(map[string]any)
	if !ok {
		return ""
	}
	return getValue(m, path[1:])
}
