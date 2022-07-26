package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"

func (r rep) SaveMetric(name string, value entity.Getter) {
	switch value := value.(type) {
	case entity.Counter:
		valueInt64 := int64(value)
		r.storage[name] = entity.Metrics{
			ID:    name,
			MType: value.GetType(),
			Delta: &valueInt64,
		}
	case entity.Gauge:
		valueFloat64 := float64(value)
		r.storage[name] = entity.Metrics{
			ID:    name,
			MType: value.GetType(),
			Value: &valueFloat64,
		}
	}
}
