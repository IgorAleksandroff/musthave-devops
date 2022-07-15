package metricscollection

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

//go:generate mockery --name Repository

type Repository interface {
	SaveGaugeMetric(name string, value float64)
	SaveCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	GetGaugeMetrics() map[string]entity.MetricGauge
	GetCounterMetrics() map[string]entity.MetricCounter
}
