package repository

import (
	"fmt"
)

func (r rep) GetGaugeMetric(name string) (float64, error) {
	if metric, ok := r.gaugeDB[name]; ok {
		return metric.Value, nil
	}

	return 0, fmt.Errorf("can not get a metric: %s", name)
}
