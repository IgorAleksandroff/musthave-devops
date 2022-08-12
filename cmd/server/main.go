package main

import (
	"context"
	"log"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/utils/serverconfig"
)

func main() {
	ctx, closeCtx := context.WithCancel(context.Background())
	defer closeCtx()

	config := serverconfig.Read()

	metricsRepo := metricscollection.NewRepository(ctx, metricscollection.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	metricsUC := metricscollection.NewUsecase(metricsRepo)

	server := api.New(config.Host, metricsUC)

	metricsRepo.MemSync()

	log.Fatal(server.Run())
}
