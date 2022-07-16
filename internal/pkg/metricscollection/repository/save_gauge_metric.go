package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (r rep) SaveGaugeMetric(name string, value float64) {
	r.gaugeDB[name] = entity.MetricGauge{
		TypeMetric: entity.GaugeTypeMetric,
		Value:      value,
	}
}
