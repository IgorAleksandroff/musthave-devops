package runtimemetrics

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(name string, value entity.Getter)
	GetMetric(name string) (m entity.Metrics, err error)
	GetMetricsName() []string
}
