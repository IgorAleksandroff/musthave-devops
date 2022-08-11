package metrichandler

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
)

type handler struct {
	metricsUC metricscollection.Usecase
}

func New(
	metricsUC metricscollection.Usecase,
) *handler {
	return &handler{
		metricsUC: metricsUC,
	}
}
