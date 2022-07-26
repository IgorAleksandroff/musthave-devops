package metrichandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (h *handler) HandleJSONGet(w http.ResponseWriter, r *http.Request) {
	var err error

	metric := Metrics{}
	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	contentTypeHeaderValue := r.Header.Get("Content-Type")
	if !strings.Contains(contentTypeHeaderValue, "application/json") {
		http.Error(w, "unknown content-type", http.StatusNotImplemented)
		return
	}

	reader := json.NewDecoder(r.Body)
	if err = reader.Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case entity.CounterTypeMetric:
		var valueMetric int64
		valueMetric, err = h.metricsUC.GetCounterMetric(metric.ID)
		metric.Delta = &valueMetric

	case entity.GaugeTypeMetric:
		var valueMetric float64
		valueMetric, err = h.metricsUC.GetGaugeMetric(metric.ID)
		metric.Value = &valueMetric

	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	err = jsonEncoder.Encode(metric)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
