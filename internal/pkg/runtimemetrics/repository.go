package runtimemetrics

import (
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/utils"
)

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(name string, value Getter)
	GetMetric(name string) (m Metrics, err error)
	GetMetricsName() []string
}

type rep struct {
	storage map[string]Metrics
	hashKey string
}

func NewRepository(key string) *rep {
	return &rep{storage: make(map[string]Metrics), hashKey: key}
}

func (r *rep) SaveMetric(name string, value Getter) {
	switch value := value.(type) {
	case Counter:
		valueInt64 := int64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Delta: &valueInt64,
			Hash:  utils.GetHash(fmt.Sprintf("%s:counter:%d", name, valueInt64), r.hashKey),
		}
	case Gauge:
		valueFloat64 := float64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Value: &valueFloat64,
			Hash:  utils.GetHash(fmt.Sprintf("%s:gauge:%f", name, valueFloat64), r.hashKey),
		}
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
