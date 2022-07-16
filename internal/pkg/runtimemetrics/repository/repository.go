package repository

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"
)

type rep struct {
	storage map[string]entity.Metric
}

func New() *rep {
	return &rep{storage: make(map[string]entity.Metric)}
}
