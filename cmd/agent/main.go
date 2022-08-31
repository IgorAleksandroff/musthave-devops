package main

import (
	"time"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	runtimemetrics2 "github.com/IgorAleksandroff/musthave-devops/internal/runtimemetrics"
)

func main() {
	config := enviroment.NewClientConfig()

	client := devopsserver.NewClient(config.Host)
	runtimeMetricsRepo := runtimemetrics2.NewRepository(config.HashKey)
	runtimeMetricsUC := runtimemetrics2.NewRuntimeMetrics(runtimeMetricsRepo, client)

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
