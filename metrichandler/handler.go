package metrichandler

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
)

type handler struct {
	metricsUC metricscollection.MetricsCollection
	hashKey   string
}

func New(
	metricsUC metricscollection.MetricsCollection,
	key string,
) *handler {
	return &handler{
		metricsUC: metricsUC,
		hashKey:   key,
	}
}
