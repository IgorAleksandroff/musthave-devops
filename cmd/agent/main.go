package main

import (
	"time"

	"github.com/IgorAleksandroff/musthave-devops/configuration/clientconfig"
	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics"
)

func main() {
	config := clientconfig.Read()

	client := devopsserver.NewClient(config.Host)
	runtimeMetricsRepo := runtimemetrics.NewRepository()
	runtimeMetricsUC := runtimemetrics.NewUsecase(runtimeMetricsRepo, client)

	pollTicker := time.NewTicker(config.PollInterval)
	reportTicker := time.NewTicker(config.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			runtimeMetricsUC.UpdateMetrics()
		case <-reportTicker.C:
			runtimeMetricsUC.SendMetrics()
		}
	}
}
