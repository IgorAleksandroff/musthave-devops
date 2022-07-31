package metricscollection

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

//go:generate mockery --name Usecase

type Usecase interface {
	SaveMetric(value entity.Metrics)
	SaveCounterMetric(value entity.Metrics)
	GetMetric(name string) (*entity.Metrics, error)
	GetMetricsValue() map[string]string
}
