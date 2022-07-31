package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/repository"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/usecase"
)

func main() {
	pollInterval := getEnvInt("POLL_INTERVAL", 2)
	reportInterval := getEnvInt("REPORT_INTERVAL", 10)
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	client := devopsserver.NewClient()
	runtimeMetricsRepo := repository.New()
	runtimeMetricsUC := usecase.New(runtimeMetricsRepo, client)

	for {
		select {
		case <-pollTicker.C:
			runtimeMetricsUC.UpdateMetrics()
		case <-reportTicker.C:
			runtimeMetricsUC.SendMetrics()
		}
	}
}

func getEnvInt(envName string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(envName))
	if err != nil {
		log.Printf("error of env %s: %s", envName, err.Error())
		return defaultValue
	}
	return value
}
