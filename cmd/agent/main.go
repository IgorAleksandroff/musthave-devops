package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/runtimemetrics"
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

	cfg := enviroment.NewClientConfig()

	client, err := devopsserver.NewClient(cfg.Host, cfg.CryptoKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	runtimeMetricsRepo := runtimemetrics.NewRepository(cfg.HashKey)
	runtimeMetricsUC := runtimemetrics.NewRuntimeMetrics(runtimeMetricsRepo, client)

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			go runtimeMetricsUC.UpdateUtilMetrics()
			go runtimeMetricsUC.UpdateMetrics()
		case <-reportTicker.C:
			go runtimeMetricsUC.SendMetricsBatch()
		}
	}
}
