// Package runtimemetrics collects runtime metrics and sends them to server via HTTP protocol.
package runtimemetrics

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
)

//go:generate mockery --name RuntimeMetrics

type RuntimeMetrics interface {
	UpdateMetrics()
	SendMetrics()
	SendMetricsBatch()
}

type runtimeMetrics struct {
	repository         Repository
	devopsServerClient devopsserver.Client
}

func NewRuntimeMetrics(
	r Repository,
	client devopsserver.Client,
) *runtimeMetrics {
	r.SaveMetric("PollCount", Counter(0))

	return &runtimeMetrics{
		repository:         r,
		devopsServerClient: client,
	}
}

// UpdateMetrics collects Memory statistics.
func (u runtimeMetrics) UpdateMetrics() {
	pollCount, err := u.repository.GetMetric("PollCount")
	if err != nil {
		log.Println(err)
		var value int64
		pollCount = Metrics{
			Delta: &value,
		}
	}
	pollCountInt := *pollCount.Delta + 1
	randomValue := float64(rand.Int())

	u.repository.SaveMetric("PollCount", Counter(pollCountInt))
	u.repository.SaveMetric("RandomValue", Gauge(randomValue))

	memMetrics := runtime.MemStats{}
	runtime.ReadMemStats(&memMetrics)

	u.repository.SaveMetric("Alloc", Gauge(float64(memMetrics.Alloc)))
	u.repository.SaveMetric("BuckHashSys", Gauge(float64(memMetrics.BuckHashSys)))
	u.repository.SaveMetric("Frees", Gauge(float64(memMetrics.Frees)))
	u.repository.SaveMetric("GCCPUFraction", Gauge(memMetrics.GCCPUFraction))
	u.repository.SaveMetric("GCSys", Gauge(float64(memMetrics.GCSys)))
	u.repository.SaveMetric("HeapAlloc", Gauge(float64(memMetrics.HeapAlloc)))
	u.repository.SaveMetric("HeapIdle", Gauge(float64(memMetrics.HeapIdle)))
	u.repository.SaveMetric("HeapInuse", Gauge(float64(memMetrics.HeapInuse)))
	u.repository.SaveMetric("HeapObjects", Gauge(float64(memMetrics.HeapObjects)))
	u.repository.SaveMetric("HeapReleased", Gauge(float64(memMetrics.HeapReleased)))
	u.repository.SaveMetric("HeapSys", Gauge(float64(memMetrics.HeapSys)))
	u.repository.SaveMetric("LastGC", Gauge(float64(memMetrics.LastGC)))
	u.repository.SaveMetric("Lookups", Gauge(float64(memMetrics.Lookups)))
	u.repository.SaveMetric("MCacheInuse", Gauge(float64(memMetrics.MCacheInuse)))
	u.repository.SaveMetric("MCacheSys", Gauge(float64(memMetrics.MCacheSys)))
	u.repository.SaveMetric("MSpanInuse", Gauge(float64(memMetrics.MSpanInuse)))
	u.repository.SaveMetric("MSpanSys", Gauge(float64(memMetrics.MSpanSys)))
	u.repository.SaveMetric("Mallocs", Gauge(float64(memMetrics.Mallocs)))
	u.repository.SaveMetric("NextGC", Gauge(float64(memMetrics.NextGC)))
	u.repository.SaveMetric("NumForcedGC", Gauge(float64(memMetrics.NumForcedGC)))
	u.repository.SaveMetric("NumGC", Gauge(float64(memMetrics.NumGC)))
	u.repository.SaveMetric("OtherSys", Gauge(float64(memMetrics.OtherSys)))
	u.repository.SaveMetric("PauseTotalNs", Gauge(float64(memMetrics.PauseTotalNs)))
	u.repository.SaveMetric("StackInuse", Gauge(float64(memMetrics.StackInuse)))
	u.repository.SaveMetric("StackSys", Gauge(float64(memMetrics.StackSys)))
	u.repository.SaveMetric("Sys", Gauge(float64(memMetrics.Sys)))
	u.repository.SaveMetric("TotalAlloc", Gauge(float64(memMetrics.TotalAlloc)))
}

// UpdateUtilMetrics collects VirtualMemory and CPU statistics.
func (u runtimeMetrics) UpdateUtilMetrics() {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return
	}

	u.repository.SaveMetric("TotalMemory", Gauge(float64(v.Total)))
	u.repository.SaveMetric("FreeMemory", Gauge(float64(v.Free)))

	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		log.Println(err)
		return
	}

	u.repository.SaveMetric("CPUutilization1", Gauge(cpuUtilization[0]))

}

// SendMetrics sends each collected metric to server via HTTP.
func (u runtimeMetrics) SendMetrics() {
	metricsName := u.repository.GetMetricsName()
	for _, metricName := range metricsName {
		m, err := u.repository.GetMetric(metricName)
		if err != nil {
			log.Println(err)
			continue
		}

		if err = u.devopsServerClient.Update(devopsserver.EndpointUpdate, []devopsserver.Metrics{
			{
				ID:    m.ID,
				MType: m.MType,
				Delta: m.Delta,
				Value: m.Value,
				Hash:  m.Hash,
			},
		}); err != nil {
			log.Println(err)
		}
	}
}

// SendMetricsBatch sends batch collected metrics to server via HTTP.
func (u runtimeMetrics) SendMetricsBatch() {
	metricsName := u.repository.GetMetrics()

	metrics := make([]devopsserver.Metrics, 0, len(metricsName))
	for _, m := range metricsName {
		metrics = append(metrics, devopsserver.Metrics{
			ID:    m.ID,
			MType: m.MType,
			Delta: m.Delta,
			Value: m.Value,
			Hash:  m.Hash,
		})
	}

	if err := u.devopsServerClient.Update(devopsserver.EndpointUpdates, metrics); err != nil {
		log.Println(err)
	}
}
