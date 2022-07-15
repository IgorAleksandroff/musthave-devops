package usecase

func (u usecase) SaveGaugeMetric(name string, value float64) {
	u.repository.SaveGaugeMetric(name, value)
}
