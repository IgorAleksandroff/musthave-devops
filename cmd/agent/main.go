package main

import (
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics"
	"github.com/IgorAleksandroff/musthave-devops/utils/clientconfig"
)

func main() {
	config := clientconfig.Read()

	client := devopsserver.NewClient(config.Host)
	runtimeMetricsRepo := runtimemetrics.NewRepository(config.HashKey)
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
