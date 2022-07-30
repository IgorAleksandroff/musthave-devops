package usecase

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (u usecase) GetMetric(name string) (*entity.Metrics, error) {
	return u.repository.GetMetric(name)
}
