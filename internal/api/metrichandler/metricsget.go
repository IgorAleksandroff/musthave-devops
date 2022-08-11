package metrichandler

import (
	"html/template"
	"net/http"
)

func (h *handler) HandleMetricsGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metricsValue := h.metricsUC.GetMetricsValue()
	tmpl.Execute(w, metricsValue)
}
