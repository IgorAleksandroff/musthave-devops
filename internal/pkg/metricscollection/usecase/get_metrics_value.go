package usecase

import "fmt"

func (u usecase) GetMetricsValue() map[string]string {
	counterMetrics := u.repository.GetCounterMetrics()
	gaugeMetrics := u.repository.GetGaugeMetrics()

	result := make(map[string]string, len(counterMetrics)+len(gaugeMetrics))
	for name, value := range counterMetrics {
		result[name] = fmt.Sprintf("%v", value.Value)
	}

	for name, value := range gaugeMetrics {
		result[name] = fmt.Sprintf("%v", value.Value)
	}

	return result
}
