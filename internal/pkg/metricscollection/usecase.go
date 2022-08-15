package metricscollection

import (
	"fmt"
	"log"
)

//go:generate mockery --name Usecase

type Usecase interface {
	SaveMetric(value Metrics)
	SaveCounterMetric(value Metrics)
	GetMetric(name string) (*Metrics, error)
	GetMetricsValue() map[string]string
}

type usecase struct {
	repository Repository
}

func NewUsecase(
	r Repository,
) *usecase {
	return &usecase{
		repository: r,
	}
}

func (u usecase) SaveMetric(value Metrics) {
	u.repository.SaveMetric(value)
}

func (u usecase) SaveCounterMetric(value Metrics) {
	if metric, err := u.repository.GetMetric(value.ID); err == nil {
		delta := *metric.Delta + *value.Delta
		value.Delta = &delta
	}

	u.repository.SaveMetric(value)
}

func (u usecase) GetMetric(name string) (*Metrics, error) {
	return u.repository.GetMetric(name)
}

func (u usecase) GetMetricsValue() map[string]string {
	metrics := u.repository.GetMetrics()

	result := make(map[string]string, len(metrics))
	for name, m := range metrics {
		switch m.MType {
		case CounterTypeMetric:
			if m.Delta != nil {
				result[name] = fmt.Sprintf("%v", *m.Delta)
			}
		case GaugeTypeMetric:
			if m.Value != nil {
				result[name] = fmt.Sprintf("%v", *m.Value)
			}
		default:
			log.Printf("wrong type: %v of metric: %v", m.MType, m.ID)
		}
	}

	return result
}
