package repository

import "github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"

type rep struct {
	metricDB map[string]entity.Metrics
}

func New() *rep {
	return &rep{
		metricDB: make(map[string]entity.Metrics),
	}
}
