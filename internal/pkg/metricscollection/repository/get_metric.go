package repository

import (
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (r rep) GetMetric(name string) (*entity.Metrics, error) {
	if metric, ok := r.metricDB[name]; ok {
		return &metric, nil
	}

	return nil, fmt.Errorf("can not found a metric: %s", name)
}
