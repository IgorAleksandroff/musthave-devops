package main

import (
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/usecase"
)

func main() {
	client := devopsserver.NewClient()
	runtimeMetricsRepo := repository.New()
	runtimeMetricsUC := usecase.New(runtimeMetricsRepo, client)

	config := client.GetConfig()
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
