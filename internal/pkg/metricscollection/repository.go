package metricscollection

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(value Metrics)
	GetMetric(name string) (*Metrics, error)
	GetMetrics() map[string]Metrics
}
