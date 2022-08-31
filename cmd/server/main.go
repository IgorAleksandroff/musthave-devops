package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/enviroment/serverconfig"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCtx()

	config := serverconfig.NewConfig()

	metricsUC, err := metricscollection.NewUsecase(ctx, metricscollection.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
		AddressDB:     config.AddressDB,
	})
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
	defer metricsUC.Close()

	server := api.NewServer(config.Host, config.HashKey, metricsUC)

	log.Fatal(server.Run())
}
