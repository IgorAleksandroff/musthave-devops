package metrichandler

import (
	"fmt"
	"net/http"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
	"github.com/go-chi/chi"
)

func (h *handler) HandleMetricGet(w http.ResponseWriter, r *http.Request) {
	var value string
	var err error
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	metric, err := h.metricsUC.GetMetric(metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	switch metricType {
	case entity.CounterTypeMetric:
		value = fmt.Sprintf("%v", *metric.Delta)
	case entity.GaugeTypeMetric:
		value = fmt.Sprintf("%v", *metric.Value)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}
