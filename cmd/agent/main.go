package main

import (
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/usecase"
)

func main() {
	pollInterval := time.NewTicker(2 * time.Second)
	reportInterval := time.NewTicker(10 * time.Second)

	client := devopsserver.NewClient()
	runtimeMetricsRepo := repository.New()
	runtimeMetricsUC := usecase.New(runtimeMetricsRepo, client)

	for {
		select {
		case <-pollInterval.C:
			runtimeMetricsUC.UpdateMetrics()
		case <-reportInterval.C:
			runtimeMetricsUC.UpdateMetrics()
		}
	}
}
