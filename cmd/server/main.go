package main

import (
	"context"
	"log"

	"github.com/IgorAleksandroff/musthave-devops/configuration/serverconfig"
	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/usecase"
)

func main() {
	ctx, closeCtx := context.WithCancel(context.Background())
	defer closeCtx()

	config := serverconfig.Read()

	metricsRepo := repository.New(ctx, repository.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	metricsUC := usecase.New(metricsRepo)

	server := api.New(config.Host, metricsUC)

	metricsRepo.MemSync()

	log.Fatal(server.Run())
}
