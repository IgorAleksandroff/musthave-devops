package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (r rep) GetMetrics() map[string]entity.Metrics {
	result := make(map[string]entity.Metrics, len(r.metricDB))

	for name, metric := range r.metricDB {
		result[name] = entity.CopyMetric(metric)
	}

	return result
}
