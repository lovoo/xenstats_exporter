package main

import (
	"log"
	"time"

	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter implements the prometheus.Collector interface. It exposes the metrics
// of a ipmi node.
type Exporter struct {
	config       Config
	metrics      []*prometheus.GaugeVec
	duration     prometheus.Gauge
	totalScrapes prometheus.Counter
	namespace    string
	replacer     *strings.Replacer
}

type metric struct {
	metricsname string
	unit        string
	value       float64
}

// Config
type Config struct {
	Xenhost string

	Credentials struct {
		Username string
		Password string
	}
}

// NewExporter instantiates a new ipmi Exporter.
func NewExporter(config Config) *Exporter {
	e := Exporter{
		config:    config,
		namespace: "xenstats",
		duration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: *namespace,
			Name:      "exporter_last_scrape_duration_seconds",
			Help:      "The last scrape duration.",
		}),
	}

	e.metrics = []*prometheus.GaugeVec{}

	e.collect()
	return &e
}

func executeCommand(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1]).Output()
	return string(out), err
}

// Describe Describes all the registered stats metrics from the xen master.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}

	ch <- e.duration.Desc()
}

// Collect collects all the registered stats metrics from the xen master.
func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	e.collect()

	for _, m := range e.metrics {
		m.Collect(metrics)
	}

	metrics <- e.duration
}

func (e *Exporter) collect() {
	now := time.Now().UnixNano()
	var err error

	stats := NewXenstats(e.config, e.namespace)
	stats.GetDriver()

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

	e.duration.Set(float64(time.Now().UnixNano()-now) / 1000000000)

}
