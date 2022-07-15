package usecase

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/runtimemetrics/entity"
)

func (u usecase) UpdateMetrics() {
	pollCount, err := u.repository.GetMetric("PollCount")
	if err != nil {
		log.Println(err)
		pollCount = entity.Metric{Value: entity.Counter(0)}
	}
	pollCountInt := int64(pollCount.Value.(entity.Counter)) + 1
	randomValue := int64(rand.Int())

	u.repository.SaveMetric("PollCount", entity.Counter(pollCountInt))
	u.repository.SaveMetric("RandomValue", entity.Counter(randomValue))

	memMetrics := runtime.MemStats{}
	runtime.ReadMemStats(&memMetrics)

	u.repository.SaveMetric("Alloc", entity.Gauge(float64(memMetrics.Alloc)))
	u.repository.SaveMetric("BuckHashSys", entity.Gauge(float64(memMetrics.BuckHashSys)))
	u.repository.SaveMetric("Frees", entity.Gauge(float64(memMetrics.Frees)))
	u.repository.SaveMetric("GCCPUFraction", entity.Gauge(memMetrics.GCCPUFraction))
	u.repository.SaveMetric("GCSys", entity.Gauge(float64(memMetrics.GCSys)))
	u.repository.SaveMetric("HeapAlloc", entity.Gauge(float64(memMetrics.HeapAlloc)))
	u.repository.SaveMetric("HeapIdle", entity.Gauge(float64(memMetrics.HeapIdle)))
	u.repository.SaveMetric("HeapInuse", entity.Gauge(float64(memMetrics.HeapInuse)))
	u.repository.SaveMetric("HeapObjects", entity.Gauge(float64(memMetrics.HeapObjects)))
	u.repository.SaveMetric("HeapReleased", entity.Gauge(float64(memMetrics.HeapReleased)))
	u.repository.SaveMetric("HeapSys", entity.Gauge(float64(memMetrics.HeapSys)))
	u.repository.SaveMetric("LastGC", entity.Gauge(float64(memMetrics.LastGC)))
	u.repository.SaveMetric("Lookups", entity.Gauge(float64(memMetrics.Lookups)))
	u.repository.SaveMetric("MCacheInuse", entity.Gauge(float64(memMetrics.MCacheInuse)))
	u.repository.SaveMetric("MCacheSys", entity.Gauge(float64(memMetrics.MCacheSys)))
	u.repository.SaveMetric("MSpanInuse", entity.Gauge(float64(memMetrics.MSpanInuse)))
	u.repository.SaveMetric("MSpanSys", entity.Gauge(float64(memMetrics.MSpanSys)))
	u.repository.SaveMetric("Mallocs", entity.Gauge(float64(memMetrics.Mallocs)))
	u.repository.SaveMetric("NextGC", entity.Gauge(float64(memMetrics.NextGC)))
	u.repository.SaveMetric("NumForcedGC", entity.Gauge(float64(memMetrics.NumForcedGC)))
	u.repository.SaveMetric("NumGC", entity.Gauge(float64(memMetrics.NumGC)))
	u.repository.SaveMetric("OtherSys", entity.Gauge(float64(memMetrics.OtherSys)))
	u.repository.SaveMetric("PauseTotalNs", entity.Gauge(float64(memMetrics.PauseTotalNs)))
	u.repository.SaveMetric("StackInuse", entity.Gauge(float64(memMetrics.StackInuse)))
	u.repository.SaveMetric("StackSys", entity.Gauge(float64(memMetrics.StackSys)))
	u.repository.SaveMetric("Sys", entity.Gauge(float64(memMetrics.Sys)))
	u.repository.SaveMetric("TotalAlloc", entity.Gauge(float64(memMetrics.TotalAlloc)))
}
