package repository

import (
	"fmt"
)

func (r rep) GetCounterMetric(name string) (int64, error) {
	if metric, ok := r.counterDB[name]; ok {
		return metric.Value, nil
	}

	return 0, fmt.Errorf("can not get a metric: %s", name)
}
