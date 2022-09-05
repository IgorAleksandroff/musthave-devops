package metricscolllectionrepo

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
)

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(value metricscollectionentity.Metrics)
	GetMetric(name string) (*metricscollectionentity.Metrics, error)
	GetMetrics() map[string]metricscollectionentity.Metrics
	Ping() error
	Close()
	MemSync()
}
