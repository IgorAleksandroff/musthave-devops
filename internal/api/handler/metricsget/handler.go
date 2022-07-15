package metricsget

import (
	"html/template"
	"net/http"

	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/metricscollection"
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
	tmpl, err := template.ParseFiles("templates/Metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metricsValue := h.metricsUC.GetMetricsValue()
	tmpl.Execute(w, metricsValue)
}
