package runtimemetrics

import (
	"fmt"
	"sync"
)

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(name string, value Getter)
	GetMetric(name string) (m Metrics, err error)
	GetMetricsName() []string
	GetMetrics() []Metrics
}

type rep struct {
	storage map[string]Metrics
	mu      sync.Mutex
	hashKey string
}

func NewRepository(key string) *rep {
	return &rep{storage: make(map[string]Metrics), mu: sync.Mutex{}, hashKey: key}
}

func (r *rep) SaveMetric(name string, value Getter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch value := value.(type) {
	case Counter:
		valueInt64 := int64(value)
		metric := Metrics{
			ID:    name,
			MType: value.GetType(),
			Delta: &valueInt64,
		}
		metric.CalcHash(fmt.Sprintf("%s:counter:%d", name, valueInt64), r.hashKey)
		r.storage[name] = metric
	case Gauge:
		valueFloat64 := float64(value)
		metric := Metrics{
			ID:    name,
			MType: value.GetType(),
			Value: &valueFloat64,
		}
		metric.CalcHash(fmt.Sprintf("%s:gauge:%f", name, valueFloat64), r.hashKey)
		r.storage[name] = metric
	}
	//log.Println(r.storage[name])
}

func (r *rep) GetMetric(name string) (m Metrics, err error) {
	if metric, ok := r.storage[name]; ok {
		return metric, nil
	}

	return m, fmt.Errorf("can not get a metric: %s", name)
}

func (r *rep) GetMetricsName() []string {
	metricsName := make([]string, 0, len(r.storage))

	for metricName := range r.storage {
		metricsName = append(metricsName, metricName)
	}

	return metricsName
}

func (r *rep) GetMetrics() []Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()

	metrics := make([]Metrics, 0, len(r.storage))
	for _, m := range r.storage {
		metrics = append(metrics, m)
	}
	return metrics
}
