package usecase

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (u usecase) SaveMetric(value entity.Metrics) {
	u.repository.SaveMetric(value)
}
