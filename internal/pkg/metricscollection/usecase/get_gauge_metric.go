package usecase

func (u usecase) GetGaugeMetric(name string) (float64, error) {
	return u.repository.GetGaugeMetric(name)
}
