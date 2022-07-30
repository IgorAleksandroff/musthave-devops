package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

type rep struct {
	gaugeDB   map[string]entity.MetricGauge
	counterDB map[string]entity.MetricCounter
	metricDB  map[string]entity.Metrics
}

func New() *rep {
	return &rep{
		gaugeDB:   make(map[string]entity.MetricGauge),
		counterDB: make(map[string]entity.MetricCounter),
		metricDB:  make(map[string]entity.Metrics),
	}
}
