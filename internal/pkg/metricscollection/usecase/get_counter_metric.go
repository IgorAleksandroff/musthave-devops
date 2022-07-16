package usecase

func (u usecase) GetCounterMetric(name string) (int64, error) {
	return u.repository.GetCounterMetric(name)
}
