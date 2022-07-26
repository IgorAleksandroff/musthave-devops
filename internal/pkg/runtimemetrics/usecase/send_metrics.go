package usecase

import (
	"log"
)

func (u usecase) SendMetrics() {
	metricsName := u.repository.GetMetricsName()
	for _, metricName := range metricsName {
		metric, err := u.repository.GetMetric(metricName)
		if err != nil {
			log.Println(err)
			continue
		}

		endpoint := "/update/"
		if _, err = u.devopsServerClient.DoPost(endpoint, metric); err != nil {
			log.Println(err)
		}
	}
}
