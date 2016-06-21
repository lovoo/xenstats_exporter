package main

import (
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter implements the prometheus.Collector interface. It exposes the metrics
// of a ipmi node.
type Exporter struct {
	config       Config
	metrics      []*prometheus.GaugeVec
	totalScrapes prometheus.Counter
	replacer     *strings.Replacer
}

// Config -
type Config struct {
	Xenhost string

	Credentials struct {
		Username string
		Password string
	}
}

// NewExporter instantiates a new ipmi Exporter.
func NewExporter(config Config) *Exporter {
	var e = &Exporter{
		config: config,
	}

	e.metrics = []*prometheus.GaugeVec{}

	e.collect()
	return e
}

// Describe Describes all the registered stats metrics from the xen master.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}
}

// Collect collects all the registered stats metrics from the xen master.
func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	e.collect()

	for _, m := range e.metrics {
		m.Collect(metrics)
	}

}

func (e *Exporter) collect() {
	var err error

	stats := NewXenstats(e.config)
	stats.GetApiCaller()

	e.metrics, err = stats.createHostMemMetrics()
	if err != nil {
		log.Printf("Xen api error in creating host memory metrics: %v", err)
	}
	storagemetrics, err := stats.createStorageMetrics()
	if err != nil {
		log.Printf("Xen api error in creating storage metrics: %v", err)
	}
	e.metrics = append(e.metrics, storagemetrics...)

	cpumetrics, err := stats.createHostCPUMetrics()
	if err != nil {
		log.Printf("Xen api error in creating host cpu metrics: %v", err)
	}
	e.metrics = append(e.metrics, cpumetrics...)

}
