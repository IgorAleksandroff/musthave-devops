package metrichandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (h *handler) HandleJSONPost(w http.ResponseWriter, r *http.Request) {
	metric := Metrics{}
	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	contentTypeHeaderValue := r.Header.Get("content-type")
	if !strings.Contains(contentTypeHeaderValue, "application/json") {
		http.Error(w, "unknown content-type", http.StatusNotImplemented)
		return
	}

	reader := json.NewDecoder(r.Body)
	// reader.UseNumber()
	if err := reader.Decode(&metric); err != nil {
		http.Error(w, "can't parse a JSON. internal error", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case entity.CounterTypeMetric:
		if metric.Delta == nil {
			http.Error(w, "empty delta for type counter. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveCounterMetric(metric.ID, *metric.Delta)
	case entity.GaugeTypeMetric:
		if metric.Value == nil {
			http.Error(w, "empty value for type gauge. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveGaugeMetric(metric.ID, *metric.Value)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}
