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
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeCtx()

	config := serverconfig.Read()

	repositoryMemo := repositorymemo.NewRepository(ctx, repositorymemo.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	metricsUC := metricscollection.NewUsecase(repositoryMemo)

	var conn *pgxpool.Pool
	var err error
	if config.AddressDB != "" {
		conn, err = pgxpool.Connect(ctx, config.AddressDB)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		log.Printf("connect to DB: %v", conn.Config())
		defer conn.Close()

		repositoryPG := repositorypg.NewRepository(ctx, conn)
		if err = repositoryPG.Init(); err != nil {
			log.Fatalf("Init DB Error: %v\n", err)
			os.Exit(1)
		}

		metricsUC = metricscollection.NewUsecase(repositoryPG)
	} else {
		repositoryMemo.MemSync()
	}

	connectionTester := repositorypg.NewPinger(ctx, conn)

	server := api.New(config.Host, config.HashKey, metricsUC, connectionTester)

	log.Fatal(server.Run())
}
