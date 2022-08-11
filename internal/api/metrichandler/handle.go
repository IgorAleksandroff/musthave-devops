package metrichandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/go-chi/chi"
)

func (h *handler) HandleMetricPost(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	metric := metricscollection.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case metricscollection.CounterTypeMetric:
		counterValue, err := strconv.ParseInt(chi.URLParam(r, "VALUE"), 10, 64)
		if err != nil {
			http.Error(w, "can't parse a int64. internal error", http.StatusBadRequest)
			return
		}

		metric.Delta = &counterValue
		h.metricsUC.SaveCounterMetric(metric)

	case metricscollection.GaugeTypeMetric:
		gaugeValue, err := strconv.ParseFloat(chi.URLParam(r, "VALUE"), 64)
		if err != nil {
			http.Error(w, "can't parse a float64. internal error", http.StatusBadRequest)
			return
		}
		metric.Value = &gaugeValue
		h.metricsUC.SaveMetric(metric)

	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}
	w.WriteHeader(http.StatusOK)
}

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
	case metricscollection.CounterTypeMetric:
		value = fmt.Sprintf("%v", *metric.Delta)
	case metricscollection.GaugeTypeMetric:
		value = fmt.Sprintf("%v", *metric.Value)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (h *handler) HandleMetricsGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metricsValue := h.metricsUC.GetMetricsValue()
	w.Header().Set("content-type", "text/html")
	tmpl.Execute(w, metricsValue)
}

func (h *handler) HandleJSONPost(w http.ResponseWriter, r *http.Request) {
	metric := metricscollection.Metrics{}
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
	//if err := reader.Decode(&metric); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	switch metric.MType {
	case metricscollection.CounterTypeMetric:
		if metric.Delta == nil {
			http.Error(w, "empty delta for type counter. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveCounterMetric(metric)
	case metricscollection.GaugeTypeMetric:
		if metric.Value == nil {
			http.Error(w, "empty value for type gauge. internal error", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveMetric(metric)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) HandleJSONGet(w http.ResponseWriter, r *http.Request) {
	var err error

	metric := metricscollection.Metrics{}
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
	case metricscollection.CounterTypeMetric:
	case metricscollection.GaugeTypeMetric:
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