package metricscollection

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(value entity.Metrics)
	GetMetric(name string) (*entity.Metrics, error)
	GetMetrics() map[string]entity.Metrics
	FlushMemo() error
	MemSync()
}
