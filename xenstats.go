package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	xsclient "github.com/xenserver/go-xenserver-client"
)

// Xenstats -
type Xenstats struct {
	xend      *ApiCaller
	xenclient *xsclient.XenAPIClient
}

// NewXenstats -
func NewXenstats(config Config) *Xenstats {
	p := new(Xenstats)

	xend := NewApiCaller(config.Xenhost, config.Credentials.Username, config.Credentials.Password)

	// Need Login first if it is a fresh session
	xenclient, err := xend.GetXenAPIClient()
	if err != nil {
		log.Printf("service.time call error: %v", err)
	}

	p.xenclient = xenclient
	p.xend = xend

	return p
}

// GetApiCaller -
func (s Xenstats) GetApiCaller() *ApiCaller {
	return s.xend
}

func (s Xenstats) createHostTotalMemMetric(hostmetrics string, hostname string) (metric *prometheus.GaugeVec, err error) {
	memoryTotal, err := s.xend.GetSpecificValue("host_metrics.get_memory_total", hostmetrics)
	if err != nil {
		return metric, fmt.Errorf("XEN Api Error: %v", err)
	}
	memTotalInt, err := strconv.ParseInt(memoryTotal.(string), 10, 64)
	if err != nil {
		return metric, fmt.Errorf("could not parse memoryTotal: %v", err)
	}

	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      "memory_total",
		Help:      "Total memory of the xen host",
		ConstLabels: map[string]string{
			"unit": "bytes",
		},
	}, []string{"hostname"})
	labels := prometheus.Labels{"hostname": hostname}
	metric.With(labels).Set(float64(memTotalInt))

	return metric, err
}

func (s Xenstats) createHostFreeMemMetric(hostmetrics string, hostname string) (metric *prometheus.GaugeVec, err error) {
	memoryTotal, err := s.xend.GetSpecificValue("host_metrics.get_memory_free", hostmetrics)
	if err != nil {
		return metric, fmt.Errorf("XEN Api Error: %v", err)
	}
	memTotalInt, err := strconv.ParseInt(memoryTotal.(string), 10, 64)
	if err != nil {
		return metric, fmt.Errorf("value conversation error: %v", err)
	}
	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      "memory_free",
		Help:      "Total memory of the xen host",
		ConstLabels: map[string]string{
			"unit": "bytes",
		},
	}, []string{"hostname"})

	labels := prometheus.Labels{"hostname": hostname}
	metric.With(labels).Set(float64(memTotalInt))

	return metric, err
}

func (s Xenstats) createHostMemMetrics() (metrics []*prometheus.GaugeVec, err error) {
	hosts, err := s.xenclient.GetHosts()
	if err != nil {
		return metrics, fmt.Errorf("XEN Api Error: %v", err)
	}

	for _, elem := range hosts {

		hostname, err := s.xend.GetSpecificValue("host.get_name_label", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		hostmetrics, err := s.xend.GetSpecificValue("host.get_metrics", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		totalMetric, err := s.createHostTotalMemMetric(hostmetrics.(string), hostname.(string))
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		metrics = append(metrics, totalMetric)

		freeMetric, err := s.createHostFreeMemMetric(hostmetrics.(string), hostname.(string))
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, freeMetric)
	}
	return metrics, err
}

func (s Xenstats) createStorageVirtualAllocationMetrics(storagemetrics string, labelname string, uuid string, defaultStorage bool) (metric *prometheus.GaugeVec, err error) {

	valloc, err := s.xend.GetSpecificValue("SR.get_virtual_allocation", storagemetrics)
	if err != nil {
		return metric, fmt.Errorf("XEN Api Error: %v", err)
	}
	vallocint, err := strconv.ParseInt(valloc.(string), 10, 64)
	if err != nil {
		return metric, fmt.Errorf("value conversation error: %v", err)
	}

	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      "storage_virtual_allocation",
		Help:      "Memory used by virtual instantances",
		ConstLabels: map[string]string{
			"unit":            "bytes",
			"default_storage": strconv.FormatBool(defaultStorage),
			"uuid":            uuid,
		},
	}, []string{"label"})

	labels := prometheus.Labels{"label": labelname}
	metric.With(labels).Set(float64(vallocint))

	return metric, err
}

func (s Xenstats) createStoragePhysicalUtilisationMetrics(storagemetrics string, labelname string, uuid string, defaultStorage bool) (metric *prometheus.GaugeVec, err error) {

	phyutil, err := s.xend.GetSpecificValue("SR.get_physical_utilisation", storagemetrics)
	if err != nil {
		return metric, fmt.Errorf("XEN Api Error: %v", err)
	}

	phyutilInt, err := strconv.ParseInt(phyutil.(string), 10, 64)
	if err != nil {
		return metric, fmt.Errorf("value conversation error: %v", err)
	}

	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      "storage_physical_utilisation",
		Help:      "Persistent data physical utilization",
		ConstLabels: map[string]string{
			"unit":            "bytes",
			"default_storage": strconv.FormatBool(defaultStorage),
			"uuid":            uuid,
		},
	}, []string{"label"})

	labels := prometheus.Labels{"label": labelname}
	metric.With(labels).Set(float64(phyutilInt))

	return metric, err
}

func (s Xenstats) createStoragePhysicalSizeMetrics(storagemetrics string, labelname string, uuid string, defaultStorage bool) (metric *prometheus.GaugeVec, err error) {
	phySize, err := s.xend.GetSpecificValue("SR.get_physical_size", storagemetrics)
	if err != nil {
		return metric, fmt.Errorf("XEN Api Error: %v", err)
	}

	phySizeInt, err := strconv.ParseInt(phySize.(string), 10, 64)
	if err != nil {
		return metric, fmt.Errorf("value conversation error: %v", err)
	}

	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      "storage_physical_size",
		Help:      "Persistent data physical size",
		ConstLabels: map[string]string{
			"unit":            "bytes",
			"default_storage": strconv.FormatBool(defaultStorage),
			"uuid":            uuid,
		},
	}, []string{"label"})

	labels := prometheus.Labels{"label": labelname}
	metric.With(labels).Set(float64(phySizeInt))

	return metric, err
}

func (s Xenstats) createMetric(name, help, unit, labelkey string, labelvalue string, value float64) (metric *prometheus.GaugeVec, err error) {
	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      name,
		Help:      help,
		ConstLabels: map[string]string{
			"unit": unit,
		},
	}, []string{labelkey})

	labels := prometheus.Labels{labelkey: labelvalue}
	metric.With(labels).Set(float64(value))

	return metric, err
}

func Btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func (s Xenstats) createPoolMetrics() (metrics []*prometheus.GaugeVec, err error) {
	hosts, err := s.xenclient.GetPools()
	if err != nil {
		return metrics, fmt.Errorf("XEN Api Error: %v", err)
	}

	for _, elem := range hosts {
		nameLabel, err := s.xend.GetSpecificValue("pool.get_name_label", elem.Ref)

		ha_enabled, err := s.xend.GetSpecificValue("pool.get_ha_enabled", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		haEnabledInt := Btof(ha_enabled.(bool))
		haEnabledMetric, err := s.createMetric("pool_ha_enabled", "true if HA is enabled on the pool, false otherwise", "bool", "pool", nameLabel.(string), float64(haEnabledInt))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, haEnabledMetric)

		haHostFailuresToTolerate, err := s.xend.GetSpecificValue("pool.get_ha_host_failures_to_tolerate", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		haHostFailuresToTolerateInt, err := strconv.ParseInt(haHostFailuresToTolerate.(string), 10, 64)
		if err != nil {
			return metrics, fmt.Errorf("value conversation error: %v", err)
		}
		haHostFailuresToTolerateMetric, err := s.createMetric("ha_host_failures_to_tolerate", "Number of host failures to tolerate before the Pool is declared to be overcommitted", "int", "pool", nameLabel.(string), float64(haHostFailuresToTolerateInt))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, haHostFailuresToTolerateMetric)

		haAllowOvercommit, err := s.xend.GetSpecificValue("pool.get_ha_allow_overcommit", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		haAllowOvercommitInt := Btof(haAllowOvercommit.(bool))
		haAllowOvercommitMetric, err := s.createMetric("ha_allow_overcommit", "If set to false then operations which would cause the Pool to become overcommitted will be blocked.", "bool", "pool", nameLabel.(string), float64(haAllowOvercommitInt))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, haAllowOvercommitMetric)

		haOvercommitted, err := s.xend.GetSpecificValue("pool.get_ha_overcommitted", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		haOvercommittedInt := Btof(haOvercommitted.(bool))
		haOvercommittedMetric, err := s.createMetric("ha_overcommitted", "True if the Pool is considered to be overcommitted i.e. if there exist insufficient physical resources to tolerate the configured number of host failures", "bool", "pool", nameLabel.(string), float64(haOvercommittedInt))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, haOvercommittedMetric)

		wlbEnabled, err := s.xend.GetSpecificValue("pool.get_wlb_enabled", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		wlbEnabledInt := Btof(wlbEnabled.(bool))
		wlbEnabledIntMetric, err := s.createMetric("wlb_enabled", "true if workload balancing is enabled on the pool, false otherwise", "bool", "pool", nameLabel.(string), float64(wlbEnabledInt))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, wlbEnabledIntMetric)
	}
	return metrics, err
}

func (s Xenstats) createStorageMetrics() (metrics []*prometheus.GaugeVec, err error) {
	allstorages, err := s.xend.GetMultiValues("SR.get_all")
	if err != nil {
		return metrics, fmt.Errorf("XEN Api Error: %v", err)
	}
	defaultStorage, err := s.xenclient.GetDefaultSR()
	if err != nil {
		return metrics, fmt.Errorf("XEN Api Error: %v", err)
	}

	for _, elem := range allstorages {

		storagelabel, err := s.xend.GetSpecificValue("SR.get_name_label", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}

		storeageuuid, err := s.xend.GetSpecificValue("SR.get_uuid", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}

		var defaultSt = false
		if defaultStorage.Ref == elem.Ref {
			defaultSt = true
		}

		valloc, err := s.createStorageVirtualAllocationMetrics(elem.Ref, storagelabel.(string), storeageuuid.(string), defaultSt)
		if err != nil {
			return metrics, err
		}

		metrics = append(metrics, valloc)

		physicalutil, err := s.createStoragePhysicalUtilisationMetrics(elem.Ref, storagelabel.(string), storeageuuid.(string), defaultSt)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, physicalutil)

		physicalsize, err := s.createStoragePhysicalSizeMetrics(elem.Ref, storagelabel.(string), storeageuuid.(string), defaultSt)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, physicalsize)
	}

	return metrics, err
}

func (s Xenstats) createCPUMetric(name, help, unit, hostname string, value float64) (metric *prometheus.GaugeVec, err error) {
	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *namespace,
		Name:      name,
		Help:      help,
		ConstLabels: map[string]string{
			"unit": unit,
		},
	}, []string{"hostname"})

	labels := prometheus.Labels{"hostname": hostname}
	metric.With(labels).Set(float64(value))

	return metric, err
}

func (s Xenstats) createHostCPUMetrics() (metrics []*prometheus.GaugeVec, err error) {

	hosts, err := s.xenclient.GetHosts()
	if err != nil {
		log.Printf("service.time call error: %v", err)
	}
	for _, elem := range hosts {
		usedCpus := int64(0)
		hostname, err := s.xend.GetSpecificValue("host.get_name_label", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		hostcpus, err := s.xend.GetMultiValues("host.get_host_CPUs", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		vms, err := s.xend.GetMultiValues("host.get_resident_VMs", elem.Ref)
		if err != nil {
			return metrics, fmt.Errorf("XEN Api Error: %v", err)
		}
		for _, elem2 := range vms {

			vmmetrics, err := s.xend.GetSpecificValue("VM.get_metrics", elem2.Ref)
			if err != nil {
				return metrics, fmt.Errorf("XEN Api Error: %v", err)
			}

			vmCPUCount, err := s.xend.GetSpecificValue("VM_metrics.get_VCPUs_number", vmmetrics.(string))
			if err != nil {
				return metrics, fmt.Errorf("XEN Api Error: %v", err)
			}
			vmCPUCountint, err := strconv.ParseInt(vmCPUCount.(string), 10, 64)
			if err != nil {
				return metrics, fmt.Errorf("value conversation error: %v", err)
			}

			usedCpus = usedCpus + vmCPUCountint

		}

		vmsPerHost := float64(len(vms))
		vmsPerHostMetric, err := s.createCPUMetric("vms_per_host", "Number of vm´s on the xenhost", "number", hostname.(string), vmsPerHost)
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, vmsPerHostMetric)

		cpusFree := int64(len(hostcpus)) - usedCpus
		cpuUtilPercent := 100 * usedCpus / int64(len(hostcpus))
		cpuNumMetric, err := s.createCPUMetric("cpus_host_num", "Number of cpu cores on the xenhost", "bytes", hostname.(string), float64(len(hostcpus)))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, cpuNumMetric)
		cpuUtilPercentageMetric, err := s.createCPUMetric("cpus_host_util", "Used cpu cores on the xenhost in percentage", "percentage", hostname.(string), float64(cpuUtilPercent))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, cpuUtilPercentageMetric)
		cpusUsedMetric, err := s.createCPUMetric("cpus_used", "Used cpu cores on the xenhost", "number", hostname.(string), float64(usedCpus))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, cpusUsedMetric)
		cpusFreeMetric, err := s.createCPUMetric("cpus_free", "Free cpu cores on the xenhost", "number", hostname.(string), float64(cpusFree))
		if err != nil {
			return metrics, fmt.Errorf("failure during a metric creation: %v", err)
		}
		metrics = append(metrics, cpusFreeMetric)

	}
	return metrics, err
}
