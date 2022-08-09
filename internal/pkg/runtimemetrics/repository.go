package runtimemetrics

import "fmt"

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(name string, value Getter)
	GetMetric(name string) (m Metrics, err error)
	GetMetricsName() []string
}

type rep struct {
	storage map[string]Metrics
}

func NewRepository() *rep {
	return &rep{storage: make(map[string]Metrics)}
}

func (r *rep) SaveMetric(name string, value Getter) {
	switch value := value.(type) {
	case Counter:
		valueInt64 := int64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Delta: &valueInt64,
		}
	case Gauge:
		valueFloat64 := float64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Value: &valueFloat64,
		}
	}
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
