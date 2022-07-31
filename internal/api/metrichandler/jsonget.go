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

	metric := entity.Metrics{}
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
	reader.Decode(&metric)

	switch metric.MType {
	case entity.CounterTypeMetric:
	case entity.GaugeTypeMetric:
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	m, errMetric := h.metricsUC.GetMetric(metric.ID)
	if errMetric != nil || metric.MType != m.MType {
		http.Error(w, errMetric.Error(), http.StatusNotFound)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	err = jsonEncoder.Encode(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
