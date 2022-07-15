package repository

func (r rep) GetMetricsName() []string {
	metricsName := make([]string, 0, len(r.storage))

	for metricName := range r.storage {
		metricsName = append(metricsName, metricName)
	}

	return metricsName
}
