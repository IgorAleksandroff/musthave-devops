package usecase

import (
	"fmt"
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

		endpoint := fmt.Sprintf("/update/%s/%s/%v/", metric.Value.GetType(), metricName, metric.Value)
		if _, err = u.devopsServerClient.DoPost(endpoint, nil); err != nil {
			log.Println(err)
		}
	}
}
