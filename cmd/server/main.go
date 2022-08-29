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
	ctx, closeCtx := context.WithTimeout(context.Background(), 10*time.Second)
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
	connectionTester := repositorypg.NewPinger(ctx, conn)
	if config.AddressDB != "" {
		conn, err = pgxpool.Connect(ctx, config.AddressDB)
		if err != nil {
			log.Fatalf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		log.Printf("connect to DB: %v", conn.Config())
		defer conn.Close()

		repositoryPG := repositorypg.NewRepository(ctx, conn)
		if err = repositoryPG.Init(); err != nil {
			log.Fatalf("init db error: %v", err)
			os.Exit(1)
		}

		connectionTester = repositorypg.NewPinger(ctx, conn)
		metricsUC = metricscollection.NewUsecase(repositoryPG)
	} else {
		repositoryMemo.MemSync()
	}

	server := api.New(config.Host, config.HashKey, metricsUC, connectionTester)

	log.Fatal(server.Run())
}
