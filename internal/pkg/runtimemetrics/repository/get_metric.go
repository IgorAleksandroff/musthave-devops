package repository

import (
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"
)

func (r rep) GetMetric(name string) (m entity.Metrics, err error) {
	if metric, ok := r.storage[name]; ok {
		return metric, nil
	}

	return m, fmt.Errorf("can not get a metric: %s", name)
}
