package main

import (
	"context"
	"log"
	"net/http"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/metrichandler"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/usecase"
)

func main() {
	ctx, closeCtx := context.WithCancel(context.Background())
	defer closeCtx()

	server := api.New()
	config := server.GetConfig()

	metricsRepo := repository.New(ctx, repository.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	metricsUC := usecase.New(metricsRepo)
	metricHandler := metrichandler.New(metricsUC)

	server.AddHandler(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", metricHandler.HandleMetricPost)
	server.AddHandler(http.MethodGet, "/value/{TYPE}/{NAME}", metricHandler.HandleMetricGet)
	server.AddHandler(http.MethodGet, "/", metricHandler.HandleMetricsGet)
	server.AddHandler(http.MethodPost, "/update/", metricHandler.HandleJSONPost)
	server.AddHandler(http.MethodPost, "/value/", metricHandler.HandleJSONGet)

	metricsRepo.MemSync()

	log.Fatal(server.Run())
}
