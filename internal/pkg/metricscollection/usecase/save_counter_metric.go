package usecase

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (u usecase) SaveCounterMetric(value entity.Metrics) {
	if metric, err := u.repository.GetMetric(value.ID); err == nil {
		delta := *metric.Delta + *value.Delta
		value.Delta = &delta
	}

	u.repository.SaveMetric(value)
}
