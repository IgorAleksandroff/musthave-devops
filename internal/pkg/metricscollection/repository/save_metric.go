package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (r rep) SaveMetric(value entity.Metrics) {
	r.metricDB[value.ID] = value
}
