package usecase

func (u usecase) SaveCounterMetric(name string, value int64) {
	u.repository.SaveCounterMetric(name, value)
}
