package metrichandler

import (
	"net/http"
	"strconv"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
	"github.com/go-chi/chi"
)

func (h *handler) HandleMetricPost(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	switch metricType {
	case entity.CounterTypeMetric:
		counterValue, err := strconv.ParseInt(chi.URLParam(r, "VALUE"), 10, 64)
		if err != nil {
			http.Error(w, "can't parse a int64. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveCounterMetric(metricName, counterValue)
	case entity.GaugeTypeMetric:
		gaugeValue, err := strconv.ParseFloat(chi.URLParam(r, "VALUE"), 64)
		if err != nil {
			http.Error(w, "can't parse a float64. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveGaugeMetric(metricName, gaugeValue)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}
