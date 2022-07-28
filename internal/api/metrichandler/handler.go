package metrichandler

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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
