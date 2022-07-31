package usecase

import (
	"fmt"
	"log"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (u usecase) GetMetricsValue() map[string]string {
	metrics := u.repository.GetMetrics()

	result := make(map[string]string, len(metrics))
	for name, m := range metrics {
		switch m.MType {
		case entity.CounterTypeMetric:
			if m.Delta != nil {
				result[name] = fmt.Sprintf("%v", *m.Delta)
			}
		case entity.GaugeTypeMetric:
			if m.Value != nil {
				result[name] = fmt.Sprintf("%v", *m.Value)
			}
		default:
			log.Printf("wrong type: %v of metric: %v", m.MType, m.ID)
		}
	}

	return result
}
