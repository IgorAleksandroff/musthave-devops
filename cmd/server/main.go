package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
)

var (
	defaultBuildValue = "N/A"
	buildVersion      = defaultBuildValue
	buildDate         = defaultBuildValue
	buildCommit       = defaultBuildValue
)

func main() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

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
