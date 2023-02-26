package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	client, err := devopsserver.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	runtimeMetricsRepo := runtimemetrics.NewRepository(cfg.HashKey)
	runtimeMetricsUC := runtimemetrics.NewRuntimeMetrics(runtimeMetricsRepo, client)

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	stopCommand := false
	for {
		select {
		case <-pollTicker.C:
			runtimeMetricsUC.UpdateUtilMetrics()
			runtimeMetricsUC.UpdateMetrics()
		case <-reportTicker.C:
			runtimeMetricsUC.SendMetricsBatch()
		case s := <-interrupt:
			stopCommand = true
			log.Printf("got signal: %s", s)
		}

		if stopCommand {
			break
		}
	}

	log.Println("app interrupted by sys signal")
}
