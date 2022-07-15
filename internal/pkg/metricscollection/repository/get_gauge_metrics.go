package repository

import (
	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/metricscollection/entity"
)

func (r rep) GetGaugeMetrics() map[string]entity.MetricGauge {
	result := make(map[string]entity.MetricGauge, len(r.gaugeDB))

	for name, value := range r.gaugeDB {
		result[name] = value
	}

	return result
}
