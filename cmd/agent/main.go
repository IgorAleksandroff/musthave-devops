package main

import (
	"fmt"
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

	config := enviroment.NewClientConfig()

	client := devopsserver.NewClient(config.Host)
	runtimeMetricsRepo := runtimemetrics.NewRepository(config.HashKey)
	runtimeMetricsUC := runtimemetrics.NewRuntimeMetrics(runtimeMetricsRepo, client)

	pollTicker := time.NewTicker(config.PollInterval)
	reportTicker := time.NewTicker(config.ReportInterval)
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
