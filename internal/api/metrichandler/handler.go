package metrichandler

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
)

type handler struct {
	metricsUC metricscollection.Usecase
	hashKey   string
}

func New(
	metricsUC metricscollection.Usecase,
	key string,
) *handler {
	return &handler{
		metricsUC: metricsUC,
		hashKey:   key,
	}
}
