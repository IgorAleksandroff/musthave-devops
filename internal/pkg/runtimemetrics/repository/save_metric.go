package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"

func (r rep) SaveMetric(name string, value entity.Getter) {
	r.storage[name] = entity.Metric{
		Value: value,
	}
}
