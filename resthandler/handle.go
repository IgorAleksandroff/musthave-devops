// Package metrichandler stores http handlers.
// @title Monitoring API
// @description Service for saving metrics and providing read access to them
// @Version 1.0

package resthandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
)

// HandleMetricPost Saves any metrics in repository.
// @Tags Metrics
// @Summary Save metric
// @Description Saves any metrics
// @Produce text/plain
// @Param  TYPE path string true "metric type" Enums(counter,gauge)
// @Param  NAME path string true "metric id"
// @Param  VALUE path string true "metric value"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 400 {string} string "Not Implemented"
// @Router /update/{TYPE}/{NAME}/{VALUE} [post]
func (h *handler) HandleMetricPost(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	metric := metricscollectionentity.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case metricscollectionentity.CounterTypeMetric:
		counterValue, err := strconv.ParseInt(chi.URLParam(r, "VALUE"), 10, 64)
		if err != nil {
			http.Error(w, "can't parse a int64. internal error", http.StatusBadRequest)
			return
		}

		metric.Delta = &counterValue
		h.metricsUC.SaveCounterMetric(metric)

	case metricscollectionentity.GaugeTypeMetric:
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

// HandleMetricGet provides read access to metric value.
// @Tags Metrics
// @Summary Get metric
// @Description Return metric value
// @Produce text/plain
// @Param  TYPE path string true "metric type" Enums(counter,gauge)
// @Param  NAME path string true "metric id"
// @Success 200 {string} string "OK"
// @Failure 404 {string} string "Not Found"
// @Failure 400 {string} string "Not Implemented"
// @Router /value/{TYPE}/{NAME} [get]
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
	case metricscollectionentity.CounterTypeMetric:
		value = fmt.Sprintf("%v", *metric.Delta)
	case metricscollectionentity.GaugeTypeMetric:
		value = fmt.Sprintf("%v", *metric.Value)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// HandleMetricsGet provides read access to all metrics value.
// @Tags Metrics
// @Summary Get metrics
// @Description Return all metrics value
// @Produce text/html
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Internal Server Error"
// @Router / [get]
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

// HandleJSONPost Saves metric in repository.
// @Tags Metrics
// @Summary Saves metric
// @Description Saves metric in repository
// @Accept  application/json
// @Param metric body metricscollectionentity.Metrics true "Metric to save"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 400 {string} string "Not Implemented"
// @Router /update/ [post]
func (h *handler) HandleJSONPost(w http.ResponseWriter, r *http.Request) {
	metric := metricscollectionentity.Metrics{}
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

	verificationMetric := metricscollectionentity.Metrics{}
	switch metric.MType {
	case metricscollectionentity.CounterTypeMetric:

		if metric.Delta == nil {
			http.Error(w, "empty delta for type counter. internal error", http.StatusBadRequest)
			return
		}
		verificationMetric.CalcHash(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), h.hashKey)
		if h.hashKey != enviroment.ClientDefaultEnvHashKey && verificationMetric.Hash != metric.Hash {
			log.Println("hash isn't valid:", verificationMetric.Hash, metric)
			http.Error(w, "hash isn't valid", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveCounterMetric(metric)
	case metricscollectionentity.GaugeTypeMetric:
		if metric.Value == nil {
			http.Error(w, "empty value for type gauge. internal error", http.StatusBadRequest)
			return
		}
		verificationMetric.CalcHash(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), h.hashKey)
		if h.hashKey != enviroment.ClientDefaultEnvHashKey && verificationMetric.Hash != metric.Hash {
			log.Println("hash isn't valid:", verificationMetric.Hash, metric)
			http.Error(w, "hash isn't valid", http.StatusBadRequest)
			return
		}

		h.metricsUC.SaveMetric(metric)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleJSONGet provides read access to metric.
// @Tags Metrics
// @Summary Get metric
// @Description Return metric value
// @Accept  application/json
// @Produce  application/json
// @Param metrics body metricscollectionentity.Metrics true "Get Metric"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 400 {string} string "Not Implemented"
// @Router /value/ [post]
func (h *handler) HandleJSONGet(w http.ResponseWriter, r *http.Request) {
	var err error

	metric := metricscollectionentity.Metrics{}
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

	m, errMetric := h.metricsUC.GetMetric(metric.ID)
	if errMetric != nil || metric.MType != m.MType {
		http.Error(w, errMetric.Error(), http.StatusNotFound)
		return
	}
	switch metric.MType {
	case metricscollectionentity.CounterTypeMetric:
		m.CalcHash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), h.hashKey)
	case metricscollectionentity.GaugeTypeMetric:
		m.CalcHash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), h.hashKey)
	default:
		http.Error(w, "unknown handler", http.StatusNotImplemented)
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

// HandleDBPing Ping repository.
// @Tags Info
// @Summary Ping repository
// @Description Checking connection to repository
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Internal Server Error"
// @Router /ping [get]
func (h *handler) HandleDBPing(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleDBPing")
	if err := h.metricsUC.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleJSONPostBatch Saves batch metrics in repository.
// @Tags Metrics
// @Summary Saves metrics
// @Description Saves batch metrics in repository
// @Accept  application/json
// @Param metrics body metricscollectionentity.Metrics true "List of metrics to save"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 400 {string} string "Not Implemented"
// @Router /updates/ [post]
func (h *handler) HandleJSONPostBatch(w http.ResponseWriter, r *http.Request) {
	metrics := make([]metricscollectionentity.Metrics, 0)
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
	reader.Decode(&metrics)
	//if err := reader.Decode(&metric); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	h.metricsUC.SaveMetrics(metrics)

	w.WriteHeader(http.StatusOK)
}
