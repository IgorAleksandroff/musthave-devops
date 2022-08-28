package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/enviroment/serverconfig"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorymemo"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorypg"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCtx()

	config := serverconfig.NewConfig()

	repositoryMemo := repositorymemo.NewRepository(ctx, repositorymemo.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	metricsUC := metricscollection.NewUsecase(repositoryMemo)
	connectionTester := repositorypg.NewPinger(ctx)

	if config.AddressDB != "" {
		repositoryPG, err := repositorypg.NewRepository(ctx, config.AddressDB)
		if err != nil {
			log.Fatalf(err.Error())
			os.Exit(1)
		}

		metricsUC = metricscollection.NewUsecase(repositoryPG)
		connectionTester = repositoryPG
	} else {
		repositoryMemo.MemSync()
	}

	server := api.NewServer(config.Host, config.HashKey, metricsUC, connectionTester)

	log.Fatal(server.Run())
}
