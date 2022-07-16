package main

import (
	"log"
	"net/http"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/handler/metricget"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/handler/metricpost"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/handler/metricsget"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/usecase"
)

func main() {
	metricsRepo := repository.New()
	metricsUC := usecase.New(metricsRepo)

	postHandler := metricpost.New(metricsUC)
	getHandler := metricget.New(metricsUC)
	getMetricsHandler := metricsget.New(metricsUC)

	server := api.New()
	server.AddHandler(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", postHandler)
	server.AddHandler(http.MethodGet, "/value/{TYPE}/{NAME}", getHandler)
	server.AddHandler(http.MethodGet, "/", getMetricsHandler)

	log.Fatal(server.Run())
}
