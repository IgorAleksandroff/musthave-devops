package metricget

import (
	"fmt"
	"net/http"

	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/metricscollection/entity"
	"github.com/go-chi/chi"
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

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	var value string
	var err error
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	switch metricType {
	case entity.CounterTypeMetric:
		valueMetric, errMetric := h.metricsUC.GetCounterMetric(metricName)
		value = fmt.Sprintf("%v", valueMetric)
		err = errMetric

	case entity.GaugeTypeMetric:
		valueMetric, errMetric := h.metricsUC.GetGaugeMetric(metricName)
		value = fmt.Sprintf("%v", valueMetric)
		err = errMetric

	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}
