package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorymemo"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorypg"
	"github.com/IgorAleksandroff/musthave-devops/utils/enviroment/serverconfig"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCtx()

	config := serverconfig.Read()

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
		defer repositoryPG.Close()

		connectionTester = repositoryPG
		metricsUC = metricscollection.NewUsecase(repositoryPG)
	} else {
		repositoryMemo.MemSync()
	}

	server := api.New(config.Host, config.HashKey, metricsUC, connectionTester)

	log.Fatal(server.Run())
}
