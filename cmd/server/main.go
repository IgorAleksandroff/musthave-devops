package main

import (
	"log"
	"net/http"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/metrichandler"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/usecase"
)

func main() {
	metricsRepo := repository.New()
	metricsUC := usecase.New(metricsRepo)

	metricHandler := metrichandler.New(metricsUC)

	server := api.New()
	server.AddHandler(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", metricHandler.HandleMetricPost)
	server.AddHandler(http.MethodGet, "/value/{TYPE}/{NAME}", metricHandler.HandleMetricGet)
	server.AddHandler(http.MethodGet, "/", metricHandler.HandleMetricsGet)
	server.AddHandler(http.MethodPost, "/update/", metricHandler.HandleJSONPost)
	server.AddHandler(http.MethodGet, "/value/", metricHandler.HandleJSONGet)

	log.Fatal(server.Run())
}
