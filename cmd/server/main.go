package main

import (
	"context"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCtx()

	config := enviroment.NewServerConfig()

	metricsUC, err := metricscollection.NewMetricsCollection(ctx, metricscollection.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
		AddressDB:     config.AddressDB,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer metricsUC.Close()
	metricsUC.MemSync()

	server := api.NewServer(config.Host, config.HashKey, metricsUC)

	log.Fatal(server.Run())
}
