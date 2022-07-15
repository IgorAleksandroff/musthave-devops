package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

func (r rep) SaveCounterMetric(name string, value int64) {
	if storedMetric, ok := r.counterDB[name]; ok {
		value += storedMetric.Value
	}

	r.counterDB[name] = entity.MetricCounter{
		TypeMetric: entity.CounterTypeMetric,
		Value:      value,
	}
}
