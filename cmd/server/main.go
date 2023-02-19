package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	config := enviroment.NewServerConfig()

	metricsUC, err := metricscollection.NewMetricsCollection(ctx, metricscollection.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
		AddressDB:     config.AddressDB,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer metricsUC.Close()

	server, err := api.NewServer(config.ServerConfig, metricsUC)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		metricsUC.MemSync()
		wg.Done()
	}()

	server.Run(ctx)

	wg.Wait()

	log.Println("server graceful shutdown have been completed")
}
