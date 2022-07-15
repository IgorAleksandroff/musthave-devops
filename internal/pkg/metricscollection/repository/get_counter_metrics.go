package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (r rep) GetCounterMetrics() map[string]entity.MetricCounter {
	result := make(map[string]entity.MetricCounter, len(r.counterDB))

	for name, value := range r.counterDB {
		result[name] = value
	}

	return result
}
