package metrichandler

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorypg"
)

type handler struct {
	metricsUC metricscollection.Usecase
	pingDB    repositorypg.Pinger
	hashKey   string
}

func New(
	metricsUC metricscollection.Usecase,
	key string,
	pingDB repositorypg.Pinger,
) *handler {
	return &handler{
		metricsUC: metricsUC,
		hashKey:   key,
		pingDB:    pingDB,
	}
}
