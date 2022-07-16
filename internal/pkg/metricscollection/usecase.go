package metricscollection

//go:generate mockery --name Usecase

type Usecase interface {
	SaveGaugeMetric(name string, value float64)
	SaveCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	GetMetricsValue() map[string]string
}
